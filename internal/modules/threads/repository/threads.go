package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pksep/comments/internal/modules/threads/model"
)

type ThreadRepoInterface interface {
	Create(ctx context.Context) (*model.Thread, error)
}

type ThreadRepo struct {
	db *pgxpool.Pool
}

func NewThreadRepo(db *pgxpool.Pool) *ThreadRepo {
	return &ThreadRepo{db: db}
}

func (r *ThreadRepo) Create(ctx context.Context) (*model.Thread, error) {
	thread := &model.Thread{
		ID: uuid.New().String(),
	}

	_, err := r.db.Exec(ctx,
		`INSERT INTO threads (id) VALUES ($1)`,
		thread.ID,
	)
	if err != nil {
		return nil, err
	}

	return thread, nil
}
