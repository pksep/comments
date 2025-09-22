package dto

import "time"

type CommentResponse struct {
	ID              string            `json:"id"`
	EntityType      string            `json:"entity_type"`
	EntityID        string            `json:"entity_id"`
	AuthorID        string            `json:"author_id"`
	Content         string            `json:"content"`
	ParentCommentID *string           `json:"parent_comment_id,omitempty"`
	ParentAnswerID  *string           `json:"parent_answer_id,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Replies         []CommentResponse `json:"replies,omitempty"` // вложенные
}
