package exams

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/pksep/location_search_server/internal/modules/exams/api/dto"
	"github.com/pksep/location_search_server/internal/modules/exams/model"
	"github.com/pksep/location_search_server/internal/modules/exams/repository"
)

type ExamService struct {
	repo repository.ExamRepoInterface
}

func NewExamService(repo repository.ExamRepoInterface) *ExamService {
	return &ExamService{repo: repo}
}

func (s *ExamService) CreateExam(ctx *gin.Context, body dto.CreateExamDTO) (model.Exam, error) {
	exam := &model.Exam{
		ProjectID: body.ProjectID,
		UserID:    body.UserID,
		IsActive:  true,
	}

	createdExam, err := s.repo.Create(ctx.Request.Context(), exam)
	if err != nil {
		return model.Exam{}, errors.New("Ошибка: " + err.Error())
	}

	return *createdExam, nil
}

func (s *ExamService) GetExam(ctx *gin.Context, id string) (model.Exam, error) {
	exam, err := s.repo.GetByID(ctx.Request.Context(), id)
	if err != nil {
		return model.Exam{}, errors.New("Ошибка: " + err.Error())
	}
	return *exam, nil
}

func (s *ExamService) ListExams(ctx *gin.Context) ([]model.Exam, error) {
	exams, err := s.repo.List(ctx.Request.Context())
	if err != nil {
		return nil, errors.New("Ошибка: " + err.Error())
	}
	return exams, nil
}

func (s *ExamService) UpdateExam(ctx *gin.Context, id string, body dto.UpdateExamDTO) (model.Exam, error) {
	exam := &model.Exam{
		ID:        id,
		ProjectID: body.ProjectID,
		UserID:    body.UserID,
		IsActive:  body.IsActive,
	}

	updatedExam, err := s.repo.Update(ctx.Request.Context(), exam)
	if err != nil {
		return model.Exam{}, errors.New("Ошибка: " + err.Error())
	}

	return *updatedExam, nil
}

func (s *ExamService) DeleteExam(ctx *gin.Context, id string) error {
	if err := s.repo.Delete(ctx.Request.Context(), id); err != nil {
		return errors.New("Ошибка: " + err.Error())
	}
	return nil
}
