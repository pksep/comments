package dto

type CreateCommentDTO struct {
	AuthorID        string  `json:"author_id" binding:"required"`
	Content         string  `json:"content" binding:"required"`
	ThreadID         *string `json:"thread_id,omitempty"`
	AnswerCommentID  *string `json:"answer_comment_id,omitempty"`
}
