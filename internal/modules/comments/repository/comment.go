package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pksep/comments/internal/modules/comments/model"
)

// CommentRepoInterface описывает методы работы с комментариями
type CommentRepoInterface interface {
	Create(ctx context.Context, comment *model.Comment) (*model.Comment, error)
	GetByID(ctx context.Context, threadId string) (*model.Comment, error)
	Update(ctx context.Context, id string, content string, authorId string) (*model.Comment, error)
	Delete(ctx context.Context, id string, authorId string) (*model.Comment, error)
	ListWithReplies(ctx context.Context, ids []string, replyLimit int) ([]model.Comment, error)
}

// CommentRepo — реализация репозитория комментариев
type CommentRepo struct {
	db *pgxpool.Pool
}

// NewCommentRepo создаёт новый репозиторий
func NewCommentRepo(db *pgxpool.Pool) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	// 1. Ensure the comment has a ThreadID
	if comment.ThreadID == nil {
		threadID := uuid.New().String()
		// Create a new thread
		_, err := r.db.Exec(ctx, `INSERT INTO threads (id) VALUES ($1)`, threadID)
		if err != nil {
			return nil, err
		}
		comment.ThreadID = &threadID
	} else {
		// Optional: check if the thread exists to avoid foreign key violation
		var exists bool
		err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM threads WHERE id = $1)`, *comment.ThreadID).Scan(&exists)
		if err != nil {
			return nil, err
		}
		if !exists {
			// Create the thread automatically
			_, err := r.db.Exec(ctx, `INSERT INTO threads (id) VALUES ($1)`, *comment.ThreadID)
			if err != nil {
				return nil, err
			}
		}
	}

	// 2. Assign ID and timestamps for the comment
	comment.ID = uuid.New().String()
	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now

	// 3. Insert the comment
	_, err := r.db.Exec(ctx,
		`INSERT INTO comments
            (id, author_id, content, thread_id, answer_comment_id, created_at, updated_at)
         VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		comment.ID,
		comment.AuthorID,
		comment.Content,
		comment.ThreadID,
		comment.AnswerCommentID,
		comment.CreatedAt,
		comment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// GetByID возвращает комментарий по thread_id
func (r *CommentRepo) GetByID(ctx context.Context, threadID string) (*model.Comment, error) {
	query := `
        SELECT id, author_id, content, thread_id, created_at, updated_at
        FROM comments
        WHERE thread_id = $1 AND deleted_at IS NULL
        ORDER BY created_at ASC
    `
	rows, err := r.db.Query(ctx, query, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.AuthorID, &c.Content, &c.ThreadID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.Replies = []model.Comment{}
		c.RepliesCount = 0 // initialize
		comments = append(comments, c)
	}

	if len(comments) == 0 {
		return nil, nil // no comments for this thread
	}

	// First one (oldest) is root
	root := comments[0]

	// Replies are everything after the root
	if len(comments) > 1 {
		root.Replies = comments[1:]
		root.RepliesCount = len(comments) - 1
	}

	return &root, nil
}

// Update обновляет комментарий
func (r *CommentRepo) Update(ctx context.Context, id string, content string, authorId string) (*model.Comment, error) {
	// Проверяем существование комментария и авторство
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

	// Обновляем только content и сразу возвращаем полный комментарий
	updatedComment := &model.Comment{}
	err = r.db.QueryRow(ctx, `
        UPDATE comments
        SET content = $1, status = $2, updated_at = $3
        WHERE id = $4
        RETURNING id, content, author_id, status, thread_id, created_at, updated_at
    `, content, model.CommentStatusEdited, time.Now(), id).Scan(
		&updatedComment.ID,
		&updatedComment.Content,
		&updatedComment.AuthorID,
		&updatedComment.Status,
		&updatedComment.ThreadID,
		&updatedComment.CreatedAt,
		&updatedComment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return updatedComment, nil
}

// Delete удаляет комментарий и возвращает его после удаления
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

	_, err = tx.Exec(ctx, `
		UPDATE comments
		SET deleted_at = NOW(), status = 'deleted', updated_at = NOW()
		WHERE id = $1
	`, id)
	if err != nil {
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
			_, err = tx.Exec(ctx, `
				UPDATE comments
				SET deleted_at = NOW(), status = 'deleted', updated_at = NOW()
				WHERE thread_id = $1
			`, *threadID)
			if err != nil {
				return nil, err
			}
		}
	}

	var deletedComment model.Comment
	err = tx.QueryRow(ctx, `
		SELECT id, thread_id, content, author_id, status, created_at, updated_at
		FROM comments
		WHERE id = $1
	`, id).Scan(
		&deletedComment.ID,
		&deletedComment.ThreadID,
		&deletedComment.Content,
		&deletedComment.AuthorID,
		&deletedComment.Status,
		&deletedComment.CreatedAt,
		&deletedComment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	deletedComment.IsFirstComment = isFirstComment

	return &deletedComment, nil
}

func (r *CommentRepo) ListWithReplies(ctx context.Context, threadIDs []string, replyLimit int) ([]model.Comment, error) {
	if len(threadIDs) == 0 {
		return nil, nil
	}

	query := `
        SELECT id, author_id, content, thread_id, created_at, updated_at
        FROM comments
        WHERE thread_id = ANY($1) AND deleted_at IS NULL
        ORDER BY created_at ASC
    `
	rows, err := r.db.Query(ctx, query, threadIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	threadComments := make(map[string][]model.Comment)

	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.AuthorID, &c.Content, &c.ThreadID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.Replies = []model.Comment{}
		if c.ThreadID != nil {
			threadComments[*c.ThreadID] = append(threadComments[*c.ThreadID], c)
		}
	}

	var result []model.Comment

	for _, comments := range threadComments {
		if len(comments) == 0 {
			continue
		}

		root := comments[0]

		totalReplies := len(comments) - 1
		root.RepliesCount = totalReplies

		if replyLimit <= 0 || totalReplies <= 0 {
			root.Replies = []model.Comment{}
		} else {
			start := len(comments) - replyLimit
			if start < 1 {
				start = 1
			}
			root.Replies = comments[start:]
		}

		result = append(result, root)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}
