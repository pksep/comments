package model

import (
	"time"
)

// Comment is a reusable comment entity that can be attached to any domain entity
// by specifying entity type and entity id.
type Comment struct {
	ID              string     `json:"id" db:"id"`
	AuthorID        string     `json:"author_id" db:"author_id"`
	Content         string     `json:"content" db:"content"`
	ThreadID        *string    `json:"thread_id,omitempty" db:"thread_id"`
	AnswerCommentID *string    `json:"answer_comment_id,omitempty" db:"answer_comment_id"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	Replies         []Comment  `json:"replies" db:"-"` // вложенные, не маппится в БД
}
