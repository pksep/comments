package dto

type UpdateCommentDTO struct {
	ID       string `json:"id" binding:"required"`
	Content  string `json:"content" binding:"required"`
	AuthorID string `json:"author_id" binding:"required"`
}
