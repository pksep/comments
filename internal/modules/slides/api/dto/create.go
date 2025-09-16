package dto

type CreateSlideDTO struct {
	ProjectID string `form:"project_id" binding:"required"`
	Text      string `form:"text" binding:"required"`
}