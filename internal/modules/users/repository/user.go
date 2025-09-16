package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pksep/location_search_server/internal/modules/users/model"
)

type UserRepoInterface interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]model.User, error)
}

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) (*model.User, error) {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx,
		"INSERT INTO users (id, initials, created_at, updated_at) VALUES ($1, $2, $3, $4)",
		user.ID, user.Initials, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	row := r.db.QueryRow(ctx,
		"SELECT id, initials, created_at, updated_at FROM users WHERE id=$1", id,
	)

	user := &model.User{}
	if err := row.Scan(&user.ID, &user.Initials, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) (*model.User, error) {
	user.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx,
		"UPDATE users SET initials=$1, updated_at=$2 WHERE id=$3",
		user.Initials, user.UpdatedAt, user.ID,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
	return err
}

func (r *UserRepo) List(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.Query(ctx, "SELECT id, initials, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		u := model.User{}
		if err := rows.Scan(&u.ID, &u.Initials, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
