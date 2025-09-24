package repository

import (
	"context"
	"time"
    "sort"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pksep/comments/internal/modules/comments/model"
)

// CommentRepoInterface описывает методы работы с комментариями
type CommentRepoInterface interface {
	Create(ctx context.Context, comment *model.Comment) (*model.Comment, error)
	GetByID(ctx context.Context, id string) (*model.Comment, error)
	Update(ctx context.Context, comment *model.Comment) (*model.Comment, error)
	Delete(ctx context.Context, id string) error
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



// GetByID возвращает комментарий по ID
func (r *CommentRepo) GetByID(ctx context.Context, id string) (*model.Comment, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, author_id, content, thread_id, answer_comment_id, created_at, updated_at
		   FROM comments WHERE id=$1`,
		id,
	)

	c := &model.Comment{}
	if err := row.Scan(
		&c.ID, &c.AuthorID, &c.Content,
		&c.ThreadID, &c.AnswerCommentID, &c.CreatedAt, &c.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return c, nil
}

// Update обновляет комментарий
func (r *CommentRepo) Update(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	comment.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx,
		`UPDATE comments SET content=$1, updated_at=$2 WHERE id=$3`,
		comment.Content, comment.UpdatedAt, comment.ID,
	)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// Delete удаляет комментарий
func (r *CommentRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM comments WHERE id=$1`, id)
	return err
}
func (r *CommentRepo) ListWithReplies(ctx context.Context, threadIDs []string, replyLimit int) ([]model.Comment, error) {
	if len(threadIDs) == 0 {
		return nil, nil
	}

	query := `
        SELECT id, author_id, content, thread_id, created_at, updated_at
        FROM comments
        WHERE thread_id = ANY($1)
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

		root := comments[0] // oldest comment

		// Take last replyLimit comments as replies (newest ones)
		start := len(comments) - replyLimit
		if start < 1 {
			start = 1
		}
		root.Replies = comments[start:]

		result = append(result, root)
	}

	// Sort root comments by CreatedAt DESC (newest threads first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}
