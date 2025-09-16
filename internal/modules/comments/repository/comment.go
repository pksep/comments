package repository

import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/pksep/comments/internal/modules/comments/model"
)

type CommentRepoInterface interface {
    Create(ctx context.Context, comment *model.Comment) (*model.Comment, error)
    GetByID(ctx context.Context, id string) (*model.Comment, error)
    Update(ctx context.Context, comment *model.Comment) (*model.Comment, error)
    Delete(ctx context.Context, id string) error
    ListByEntity(ctx context.Context, entityType string, entityID string) ([]model.Comment, error)
}

type CommentRepo struct {
    db *pgxpool.Pool
}

func NewCommentRepo(db *pgxpool.Pool) *CommentRepo {
    return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
    comment.ID = uuid.New().String()
    comment.CreatedAt = time.Now()
    comment.UpdatedAt = time.Now()

    _, err := r.db.Exec(ctx,
        "INSERT INTO comments (id, entity_type, entity_id, author_id, content, parent_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
        comment.ID, comment.EntityType, comment.EntityID, comment.AuthorID, comment.Content, comment.ParentID, comment.CreatedAt, comment.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return comment, nil
}

func (r *CommentRepo) GetByID(ctx context.Context, id string) (*model.Comment, error) {
    row := r.db.QueryRow(ctx,
        "SELECT id, entity_type, entity_id, author_id, content, parent_id, created_at, updated_at FROM comments WHERE id=$1",
        id,
    )

    c := &model.Comment{}
    if err := row.Scan(&c.ID, &c.EntityType, &c.EntityID, &c.AuthorID, &c.Content, &c.ParentID, &c.CreatedAt, &c.UpdatedAt); err != nil {
        return nil, err
    }
    return c, nil
}

func (r *CommentRepo) Update(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
    comment.UpdatedAt = time.Now()
    _, err := r.db.Exec(ctx,
        "UPDATE comments SET content=$1, updated_at=$2 WHERE id=$3",
        comment.Content, comment.UpdatedAt, comment.ID,
    )
    if err != nil {
        return nil, err
    }
    return comment, nil
}

func (r *CommentRepo) Delete(ctx context.Context, id string) error {
    _, err := r.db.Exec(ctx, "DELETE FROM comments WHERE id=$1", id)
    return err
}

func (r *CommentRepo) ListByEntity(ctx context.Context, entityType string, entityID string) ([]model.Comment, error) {
    rows, err := r.db.Query(ctx,
        "SELECT id, entity_type, entity_id, author_id, content, parent_id, created_at, updated_at FROM comments WHERE entity_type=$1 AND entity_id=$2 ORDER BY created_at ASC",
        entityType, entityID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    out := []model.Comment{}
    for rows.Next() {
        c := model.Comment{}
        if err := rows.Scan(&c.ID, &c.EntityType, &c.EntityID, &c.AuthorID, &c.Content, &c.ParentID, &c.CreatedAt, &c.UpdatedAt); err != nil {
            return nil, err
        }
        out = append(out, c)
    }
    return out, nil
}


