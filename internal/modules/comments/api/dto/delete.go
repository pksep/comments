package dto

type DeleteCommentDTO struct {
	ID string `json:"id" binding:"required"`
}
