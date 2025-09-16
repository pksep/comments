package dto

type CreateExamDTO struct {
	ProjectID string `json:"project_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
}