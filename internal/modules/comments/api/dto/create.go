package dto

type CreateCommentDTO struct {
	EntityType      string  `json:"entity_type" binding:"required"`
	EntityID        string  `json:"entity_id" binding:"required"`
	AuthorID        string  `json:"author_id" binding:"required"`
	Content         string  `json:"content" binding:"required"`
	ParentCommentID *string `json:"parent_comment_id,omitempty"`
	ParentAnswerID  *string `json:"parent_answer_id,omitempty"`
}
