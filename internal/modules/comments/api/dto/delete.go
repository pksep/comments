package dto

type DeleteCommentDTO struct {
	ID       string `json:"id" binding:"required"`
	AuthorID string `json:"author_id" binding:"required"`
}
