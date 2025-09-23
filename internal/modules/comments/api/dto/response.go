package dto

import "time"

type CommentResponse struct {
	ID              string            `json:"id"`
	AuthorID        string            `json:"author_id"`
	Content         string            `json:"content"`
	ParentCommentID *string           `json:"parent_comment_id,omitempty"`
	AnswerCommentID  *string           `json:"answer_comment_id,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Replies         []CommentResponse `json:"replies,omitempty"` // вложенные
}
