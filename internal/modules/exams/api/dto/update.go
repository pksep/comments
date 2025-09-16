package dto

type UpdateExamDTO struct {
	ProjectID string `json:"project_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	IsActive  bool   `json:"is_active"`
}