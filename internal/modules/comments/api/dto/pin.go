package dto

type PinCommentDTO struct {
	ID       string `json:"id" binding:"required"`
	AuthorID string `json:"author_id" binding:"required"`
}
