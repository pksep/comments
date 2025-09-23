package repository

import (
	"context"
	"time"

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

// Create вставляет новый комментарий
func (r *CommentRepo) Create(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	comment.ID = uuid.New().String()
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx,
		`INSERT INTO comments
		 (id, author_id, content, parent_comment_id, answer_comment_id, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		comment.ID, comment.AuthorID,
		comment.Content, comment.ParentCommentID, comment.AnswerCommentID,
		comment.CreatedAt, comment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// GetByID возвращает комментарий по ID
func (r *CommentRepo) GetByID(ctx context.Context, id string) (*model.Comment, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, author_id, content, parent_comment_id, answer_comment_id, created_at, updated_at
		   FROM comments WHERE id=$1`,
		id,
	)

	c := &model.Comment{}
	if err := row.Scan(
		&c.ID, &c.AuthorID, &c.Content,
		&c.ParentCommentID, &c.AnswerCommentID, &c.CreatedAt, &c.UpdatedAt,
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

// ListWithReplies возвращает root-комменты с replyLimit последних реплаев для каждого
func (r *CommentRepo) ListWithReplies(ctx context.Context, ids []string, replyLimit int) ([]model.Comment, error) {
	// 1. Достаем root-комменты
	query := `
        SELECT id, author_id, content,
               parent_comment_id, answer_comment_id, created_at, updated_at
          FROM comments
         WHERE parent_comment_id IS NULL
    `
	args := []any{}
	if len(ids) > 0 {
		query += " AND id = ANY($1)"
		args = append(args, ids)
	}
	query += " ORDER BY created_at ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roots []model.Comment
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.AuthorID, &c.Content,
			&c.ParentCommentID, &c.AnswerCommentID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.Replies = []model.Comment{}
		roots = append(roots, c)
	}

	// 2. Для каждого root достаем replyLimit реплаев
	for i := range roots {
		replyRows, err := r.db.Query(ctx,
			`SELECT id, author_id, content,
                    parent_comment_id, answer_comment_id, created_at, updated_at
               FROM comments
              WHERE parent_comment_id = $1
              ORDER BY created_at DESC
              LIMIT $2`,
			roots[i].ID, replyLimit,
		)
		if err != nil {
			return nil, err
		}

		for replyRows.Next() {
			var rply model.Comment
			if err := replyRows.Scan(&rply.ID, &rply.AuthorID, &rply.Content,
				&rply.ParentCommentID, &rply.AnswerCommentID, &rply.CreatedAt, &rply.UpdatedAt); err != nil {
				replyRows.Close()
				return nil, err
			}
			rply.Replies = []model.Comment{}
			roots[i].Replies = append(roots[i].Replies, rply)
		}
		replyRows.Close()
	}

	return roots, nil
}
