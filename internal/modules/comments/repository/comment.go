package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"         // для pgx.Rows
	"github.com/jackc/pgx/v5/pgxpool" // для пула соединений
	"github.com/pksep/comments/internal/modules/comments/model"
)

type CommentRepoInterface interface {
	Create(ctx context.Context, comment *model.Comment) (*model.Comment, error)
	GetByID(ctx context.Context, id string) (*model.Comment, error)
	Update(ctx context.Context, comment *model.Comment) (*model.Comment, error)
	Delete(ctx context.Context, id string) error
	ListByEntity(ctx context.Context, entityType string, entityID string) ([]model.Comment, error)
	ListWithReplies(ctx context.Context, ids []string, depth int) ([]model.Comment, error)
}

type CommentRepo struct {
	db *pgxpool.Pool
}

func NewCommentRepo(db *pgxpool.Pool) *CommentRepo {
	return &CommentRepo{db: db}
}

// Create inserts a new comment into the database
func (r *CommentRepo) Create(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	comment.ID = uuid.New().String()
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx,
		`INSERT INTO comments
		 (id, entity_type, entity_id, author_id, content,
		  parent_comment_id, parent_answer_id, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		comment.ID, comment.EntityType, comment.EntityID, comment.AuthorID,
		comment.Content, comment.ParentCommentID, comment.ParentAnswerID,
		comment.CreatedAt, comment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// GetByID returns a comment by its ID
func (r *CommentRepo) GetByID(ctx context.Context, id string) (*model.Comment, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, entity_type, entity_id, author_id, content,
		        parent_comment_id, parent_answer_id, created_at, updated_at
		   FROM comments WHERE id=$1`,
		id,
	)

	c := &model.Comment{}
	if err := row.Scan(
		&c.ID, &c.EntityType, &c.EntityID, &c.AuthorID, &c.Content,
		&c.ParentCommentID, &c.ParentAnswerID, &c.CreatedAt, &c.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return c, nil
}

// Update updates the content of an existing comment
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

// Delete removes a comment by ID
func (r *CommentRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM comments WHERE id=$1`, id)
	return err
}

// ListByEntity returns all comments for a given entity type and ID, structured as a tree
func (r *CommentRepo) ListByEntity(ctx context.Context, entityType string, entityID string) ([]model.Comment, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, entity_type, entity_id, author_id, content,
		        parent_comment_id, parent_answer_id, created_at, updated_at
		   FROM comments
		  WHERE entity_type=$1 AND entity_id=$2
		  ORDER BY created_at ASC`,
		entityType, entityID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	commentMap := make(map[string]*model.Comment)
	all := []*model.Comment{}

	for rows.Next() {
		c := &model.Comment{}
		if err := rows.Scan(
			&c.ID, &c.EntityType, &c.EntityID, &c.AuthorID, &c.Content,
			&c.ParentCommentID, &c.ParentAnswerID, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		c.Replies = []model.Comment{}
		commentMap[c.ID] = c
		all = append(all, c)
	}

	// Build tree
	var roots []model.Comment
	for _, c := range all {
		if c.ParentCommentID != nil {
			if parent, ok := commentMap[*c.ParentCommentID]; ok {
				parent.Replies = append(parent.Replies, *c)
				continue
			}
		}
		roots = append(roots, *c)
	}

	return roots, nil
}

// ListWithReplies returns comments by IDs including nested replies up to specified depth
func (r *CommentRepo) ListWithReplies(ctx context.Context, ids []string, depth int) ([]model.Comment, error) {
	if depth <= 0 {
		return []model.Comment{}, nil
	}

	var rows pgx.Rows
	var err error

	if len(ids) == 0 {
		rows, err = r.db.Query(ctx,
			`SELECT id, entity_type, entity_id, author_id, content,
			        parent_comment_id, parent_answer_id, created_at, updated_at
			   FROM comments
			  WHERE parent_comment_id IS NULL
			  ORDER BY created_at ASC`,
		)
	} else {
		rows, err = r.db.Query(ctx,
			`SELECT id, entity_type, entity_id, author_id, content,
			        parent_comment_id, parent_answer_id, created_at, updated_at
			   FROM comments
			  WHERE parent_comment_id = ANY($1)
			  ORDER BY created_at ASC`,
			ids,
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	commentMap := make(map[string]*model.Comment)
	all := []*model.Comment{}

	for rows.Next() {
		c := &model.Comment{}
		if err := rows.Scan(
			&c.ID, &c.EntityType, &c.EntityID, &c.AuthorID, &c.Content,
			&c.ParentCommentID, &c.ParentAnswerID, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		c.Replies = []model.Comment{}
		commentMap[c.ID] = c
		all = append(all, c)
	}

	// рекурсивно подгружаем детей
	if len(all) > 0 && depth > 1 {
		childIDs := []string{}
		for _, c := range all {
			childIDs = append(childIDs, c.ID)
		}

		children, err := r.ListWithReplies(ctx, childIDs, depth-1)
		if err != nil {
			return nil, err
		}

		for _, child := range children {
			if child.ParentCommentID != nil {
				if parent, ok := commentMap[*child.ParentCommentID]; ok {
					parent.Replies = append(parent.Replies, child)
				}
			}
		}
	}

	result := make([]model.Comment, len(all))
	for i, c := range all {
		result[i] = *c
	}

	return result, nil
}
