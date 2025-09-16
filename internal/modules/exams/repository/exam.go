package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pksep/location_search_server/internal/modules/exams/model"
)

type ExamRepoInterface interface {
	Create(ctx context.Context, exam *model.Exam) (*model.Exam, error)
	GetByID(ctx context.Context, id string) (*model.Exam, error)
	Update(ctx context.Context, exam *model.Exam) (*model.Exam, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]model.Exam, error)
}

type ExamRepo struct {
	db *pgxpool.Pool
}

func NewExamRepo(db *pgxpool.Pool) *ExamRepo {
	return &ExamRepo{db: db}
}

func (r *ExamRepo) Create(ctx context.Context, exam *model.Exam) (*model.Exam, error) {
	exam.ID = uuid.New().String()
	exam.CreatedAt = time.Now()
	exam.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx,
		"INSERT INTO exams (id, project_id, user_id, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		exam.ID, exam.ProjectID, exam.UserID, exam.IsActive, exam.CreatedAt, exam.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return exam, nil
}

func (r *ExamRepo) GetByID(ctx context.Context, id string) (*model.Exam, error) {
	row := r.db.QueryRow(ctx,
		"SELECT id, project_id, user_id, is_active, created_at, updated_at FROM exams WHERE id=$1", id,
	)

	exam := &model.Exam{}
	if err := row.Scan(&exam.ID, &exam.ProjectID, &exam.UserID, &exam.IsActive, &exam.CreatedAt, &exam.UpdatedAt); err != nil {
		return nil, err
	}
	return exam, nil
}

func (r *ExamRepo) Update(ctx context.Context, exam *model.Exam) (*model.Exam, error) {
	exam.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx,
		"UPDATE exams SET project_id=$1, user_id=$2, is_active=$3, updated_at=$4 WHERE id=$5",
		exam.ProjectID, exam.UserID, exam.IsActive, exam.UpdatedAt, exam.ID,
	)
	if err != nil {
		return nil, err
	}
	return exam, nil
}

func (r *ExamRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM exams WHERE id=$1", id)
	return err
}

func (r *ExamRepo) List(ctx context.Context) ([]model.Exam, error) {
	rows, err := r.db.Query(ctx, "SELECT id, project_id, user_id, is_active, created_at, updated_at FROM exams")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exams := []model.Exam{}
	for rows.Next() {
		e := model.Exam{}
		if err := rows.Scan(&e.ID, &e.ProjectID, &e.UserID, &e.IsActive, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		exams = append(exams, e)
	}
	return exams, nil
}
