package dto

type CommentMediaDTO struct {
	ID           int64  `json:"id,omitempty"`
	Name         string `json:"name" binding:"required"`
	OriginalName string `json:"original_name,omitempty"`
	Path         string `json:"path,omitempty"`
	Type         string `json:"type" binding:"required"`
	Size         *int64 `json:"size,omitempty"`
}

type CreateCommentDTO struct {
	AuthorID        string            `json:"author_id" binding:"required"`
	Content         string            `json:"content"`
	ThreadID        *string           `json:"thread_id,omitempty"`
	AnswerCommentID *string           `json:"answer_comment_id,omitempty"`
	Documents       []CommentMediaDTO `json:"documents,omitempty"`
}
