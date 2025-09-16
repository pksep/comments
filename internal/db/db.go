package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pksep/comments/internal/config"
)

func NewPostgresPool() (*pgxpool.Pool, error) {
	cfg := config.GetConfig()
	dsn := cfg.DatabaseURL
	
	return pgxpool.New(context.Background(), dsn)
}
