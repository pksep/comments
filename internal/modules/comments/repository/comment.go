package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pksep/comments/internal/config"
	"github.com/pksep/comments/internal/modules/comments/model"
)

type CommentRepoInterface interface {
	Create(ctx context.Context, comment *model.Comment) (*model.Comment, error)
	GetByID(ctx context.Context, threadId string) (*model.Comment, error)
	Update(ctx context.Context, id string, content string, authorId string) (*model.Comment, error)
	Delete(ctx context.Context, id string, authorId string) (*model.Comment, error)
	SetPinned(ctx context.Context, id string, authorId string, isPinned bool) (*model.Comment, error)
	ListWithReplies(ctx context.Context, ids []string, replyLimit int) ([]model.Comment, error)
}

type CommentRepo struct {
	db              *pgxpool.Pool
	minioPublicURL  string
	minioBucketName string
}

type dbExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewCommentRepo(db *pgxpool.Pool) *CommentRepo {
	cfg := config.GetConfig()

	return &CommentRepo{
		db:              db,
		minioPublicURL:  strings.TrimRight(cfg.MinioPublicBaseURL, "/"),
		minioBucketName: strings.Trim(cfg.MinioBucketName, "/"),
	}
}

func (r *CommentRepo) buildMediaURL(objectName string) string {
	if objectName == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s/%s", r.minioPublicURL, r.minioBucketName, objectName)
}

func (r *CommentRepo) saveMedia(ctx context.Context, executor dbExecutor, commentID string, documents []model.CommentMedia) error {
	for index := range documents {
		document := documents[index]
		if document.Type == "" {
			document.Type = "file"
		}

		var createdAt time.Time
		var updatedAt time.Time
		err := executor.QueryRow(ctx, `
			INSERT INTO comment_media
				(comment_id, name, original_name, type, size, created_at, updated_at)
			VALUES
				($1, $2, $3, $4, $5, NOW(), NOW())
			RETURNING id, created_at, updated_at
		`,
			commentID,
			document.Name,
			document.OriginalName,
			document.Type,
			document.Size,
		).Scan(&documents[index].ID, &createdAt, &updatedAt)
		if err != nil {
			return err
		}

		documents[index].Type = document.Type
		documents[index].Path = r.buildMediaURL(document.Name)
		documents[index].CreatedAt = &createdAt
		documents[index].UpdatedAt = &updatedAt
	}

	return nil
}

func (r *CommentRepo) loadMediaMap(ctx context.Context, commentIDs []string) (map[string][]model.CommentMedia, error) {
	result := make(map[string][]model.CommentMedia)
	if len(commentIDs) == 0 {
		return result, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT comment_id, id, name, original_name, type, size, created_at, updated_at
		FROM comment_media
		WHERE comment_id = ANY($1)
		ORDER BY id ASC
	`, commentIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var commentID string
		var document model.CommentMedia
		if err := rows.Scan(
			&commentID,
			&document.ID,
			&document.Name,
			&document.OriginalName,
			&document.Type,
			&document.Size,
			&document.CreatedAt,
			&document.UpdatedAt,
		); err != nil {
			return nil, err
		}
		document.Path = r.buildMediaURL(document.Name)
		result[commentID] = append(result[commentID], document)
	}

	return result, rows.Err()
}

func collectCommentIDs(comments []model.Comment) []string {
	ids := make([]string, 0, len(comments))
	for _, comment := range comments {
		ids = append(ids, comment.ID)
	}
	return ids
}

func assignMediaToComments(comments []model.Comment, mediaMap map[string][]model.CommentMedia) []model.Comment {
	for index := range comments {
		comments[index].Documents = mediaMap[comments[index].ID]
	}
	return comments
}

func (r *CommentRepo) Create(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if comment.ThreadID == nil {
		threadID := uuid.New().String()
		if _, err := tx.Exec(ctx, `INSERT INTO threads (id) VALUES ($1)`, threadID); err != nil {
			return nil, err
		}
		comment.ThreadID = &threadID
	} else {
		var exists bool
		err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM threads WHERE id = $1)`, *comment.ThreadID).Scan(&exists)
		if err != nil {
			return nil, err
		}
		if !exists {
			if _, err := tx.Exec(ctx, `INSERT INTO threads (id) VALUES ($1)`, *comment.ThreadID); err != nil {
				return nil, err
			}
		}
	}

	comment.ID = uuid.New().String()
	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now

	_, err = tx.Exec(ctx, `
		INSERT INTO comments
			(id, author_id, content, thread_id, answer_comment_id, is_pinned, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`,
		comment.ID,
		comment.AuthorID,
		comment.Content,
		comment.ThreadID,
		comment.AnswerCommentID,
		comment.IsPinned,
		comment.CreatedAt,
		comment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := r.saveMedia(ctx, tx, comment.ID, comment.Documents); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *CommentRepo) GetByID(ctx context.Context, threadID string) (*model.Comment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, author_id, content, thread_id, answer_comment_id, is_pinned, status, created_at, updated_at
		FROM comments
		WHERE thread_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
	`, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.AuthorID, &comment.Content, &comment.ThreadID, &comment.AnswerCommentID, &comment.IsPinned, &comment.Status, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, err
		}
		comment.Replies = []model.Comment{}
		comment.Documents = []model.CommentMedia{}
		comments = append(comments, comment)
	}

	if len(comments) == 0 {
		return nil, nil
	}

	mediaMap, err := r.loadMediaMap(ctx, collectCommentIDs(comments))
	if err != nil {
		return nil, err
	}
	comments = assignMediaToComments(comments, mediaMap)

	root := comments[0]
	if len(comments) > 1 {
		root.Replies = comments[1:]
		sort.SliceStable(root.Replies, func(i, j int) bool {
			return commentComesFirst(root.Replies[i], root.Replies[j], true)
		})
		root.RepliesCount = len(comments) - 1
	}

	return &root, nil
}

func (r *CommentRepo) Update(ctx context.Context, id string, content string, authorId string) (*model.Comment, error) {
	var dbAuthor string
	err := r.db.QueryRow(ctx, `
		SELECT author_id
		FROM comments
		WHERE id = $1
	`, id).Scan(&dbAuthor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("comment with ID %s not found", id)
		}
		return nil, err
	}

	if dbAuthor != authorId {
		return nil, fmt.Errorf("only the author can edit this comment")
	}

	updatedComment := &model.Comment{}
	err = r.db.QueryRow(ctx, `
		UPDATE comments
		SET content = $1, status = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, content, author_id, is_pinned, status, thread_id, answer_comment_id, created_at, updated_at
	`, content, model.CommentStatusEdited, time.Now(), id).Scan(
		&updatedComment.ID,
		&updatedComment.Content,
		&updatedComment.AuthorID,
		&updatedComment.IsPinned,
		&updatedComment.Status,
		&updatedComment.ThreadID,
		&updatedComment.AnswerCommentID,
		&updatedComment.CreatedAt,
		&updatedComment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	mediaMap, err := r.loadMediaMap(ctx, []string{updatedComment.ID})
	if err != nil {
		return nil, err
	}
	updatedComment.Documents = mediaMap[updatedComment.ID]

	return updatedComment, nil
}

func (r *CommentRepo) Delete(ctx context.Context, id string, authorId string) (*model.Comment, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var dbAuthor string
	var threadID *string
	err = tx.QueryRow(ctx, `
		SELECT author_id, thread_id
		FROM comments
		WHERE id = $1
	`, id).Scan(&dbAuthor, &threadID)
	if err != nil {
		return nil, err
	}

	allowed := dbAuthor == authorId
	isFirstComment := false

	if !allowed && threadID != nil {
		var threadAuthor string
		err = tx.QueryRow(ctx, `
			SELECT author_id
			FROM comments
			WHERE thread_id = $1
			ORDER BY created_at ASC
			LIMIT 1
		`, *threadID).Scan(&threadAuthor)
		if err != nil {
			return nil, err
		}
		if threadAuthor == authorId {
			allowed = true
		}
	}

	if !allowed {
		return nil, errors.New("only the comment author or thread author can delete this comment")
	}

	if _, err := tx.Exec(ctx, `
		UPDATE comments
		SET deleted_at = NOW(), status = 'deleted', updated_at = NOW()
		WHERE id = $1
	`, id); err != nil {
		return nil, err
	}

	if threadID != nil {
		var firstCommentID string
		err = tx.QueryRow(ctx, `
			SELECT id FROM comments
			WHERE thread_id = $1
			ORDER BY created_at ASC
			LIMIT 1
		`, *threadID).Scan(&firstCommentID)
		if err != nil {
			return nil, err
		}

		if firstCommentID == id {
			isFirstComment = true
			if _, err := tx.Exec(ctx, `
				UPDATE comments
				SET deleted_at = NOW(), status = 'deleted', updated_at = NOW()
				WHERE thread_id = $1
			`, *threadID); err != nil {
				return nil, err
			}
		}
	}

	var deletedComment model.Comment
	err = tx.QueryRow(ctx, `
		SELECT id, thread_id, content, author_id, answer_comment_id, status, created_at, updated_at, is_pinned
		FROM comments
		WHERE id = $1
	`, id).Scan(
		&deletedComment.ID,
		&deletedComment.ThreadID,
		&deletedComment.Content,
		&deletedComment.AuthorID,
		&deletedComment.AnswerCommentID,
		&deletedComment.Status,
		&deletedComment.CreatedAt,
		&deletedComment.UpdatedAt,
		&deletedComment.IsPinned,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	mediaMap, err := r.loadMediaMap(ctx, []string{deletedComment.ID})
	if err != nil {
		return nil, err
	}
	deletedComment.Documents = mediaMap[deletedComment.ID]
	deletedComment.IsFirstComment = isFirstComment

	return &deletedComment, nil
}

func (r *CommentRepo) SetPinned(ctx context.Context, id string, authorId string, isPinned bool) (*model.Comment, error) {
	var dbAuthor string
	err := r.db.QueryRow(ctx, `
		SELECT author_id
		FROM comments
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&dbAuthor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("comment with ID %s not found", id)
		}
		return nil, err
	}

	if dbAuthor != authorId {
		return nil, fmt.Errorf("only the author can pin or unpin this comment")
	}

	item := &model.Comment{}
	err = r.db.QueryRow(ctx, `
		UPDATE comments
		SET is_pinned = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, author_id, content, thread_id, answer_comment_id, is_pinned, status, created_at, updated_at
	`, isPinned, time.Now(), id).Scan(
		&item.ID,
		&item.AuthorID,
		&item.Content,
		&item.ThreadID,
		&item.AnswerCommentID,
		&item.IsPinned,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	mediaMap, err := r.loadMediaMap(ctx, []string{item.ID})
	if err != nil {
		return nil, err
	}
	item.Documents = mediaMap[item.ID]

	return item, nil
}

func (r *CommentRepo) ListWithReplies(ctx context.Context, threadIDs []string, replyLimit int) ([]model.Comment, error) {
	if len(threadIDs) == 0 {
		return nil, nil
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, author_id, content, thread_id, answer_comment_id, is_pinned, status, created_at, updated_at
		FROM comments
		WHERE thread_id = ANY($1) AND deleted_at IS NULL
		ORDER BY created_at ASC
	`, threadIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threadComments := make(map[string][]model.Comment)
	var allComments []model.Comment

	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.AuthorID, &comment.Content, &comment.ThreadID, &comment.AnswerCommentID, &comment.IsPinned, &comment.Status, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, err
		}
		comment.Replies = []model.Comment{}
		comment.Documents = []model.CommentMedia{}
		if comment.ThreadID != nil {
			threadComments[*comment.ThreadID] = append(threadComments[*comment.ThreadID], comment)
		}
		allComments = append(allComments, comment)
	}

	mediaMap, err := r.loadMediaMap(ctx, collectCommentIDs(allComments))
	if err != nil {
		return nil, err
	}

	var result []model.Comment
	for _, comments := range threadComments {
		if len(comments) == 0 {
			continue
		}

		comments = assignMediaToComments(comments, mediaMap)
		root := comments[0]
		totalReplies := len(comments) - 1
		root.RepliesCount = totalReplies

		if replyLimit > 0 && totalReplies > 0 {
			replies := append([]model.Comment(nil), comments[1:]...)
			sort.SliceStable(replies, func(i, j int) bool {
				return commentComesFirst(replies[i], replies[j], false)
			})
			if len(replies) > replyLimit {
				replies = replies[:replyLimit]
			}
			root.Replies = replies
		}

		result = append(result, root)
	}

	sort.Slice(result, func(i, j int) bool {
		return commentComesFirst(result[i], result[j], false)
	})

	return result, nil
}
