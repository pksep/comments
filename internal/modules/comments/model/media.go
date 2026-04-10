package model

import "time"

type CommentMedia struct {
	ID           int64      `json:"id,omitempty" db:"id"`
	Name         string     `json:"name" db:"name"`
	OriginalName *string    `json:"original_name,omitempty" db:"original_name"`
	Path         string     `json:"path"`
	Type         string     `json:"type" db:"type"`
	Size         *int64     `json:"size,omitempty" db:"size"`
	CreatedAt    *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
