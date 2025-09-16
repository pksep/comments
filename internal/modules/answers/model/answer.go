package model

import (
	"github.com/pksep/location_search_server/internal/modules/shared/model"
)

// Answer represents an answer to a slide in an exam
// @Description Answer to a slide in an exam
type Answer struct {
	// ID of the answer
	// example: 123e4567-e89b-12d3-a456-426614174004
	ID string `json:"id"`

	// ExamID is the ID of the exam this answer belongs to
	// example: 123e4567-e89b-12d3-a456-426614174003
	ExamID string `json:"exam_id" gorm:"index"`

	// SlideID is the ID of the slide this answer is for
	// example: 123e4567-e89b-12d3-a456-426614174002
	SlideID string `json:"slide_id" gorm:"index"`

	// IsCorrect indicates if the answer is correct
	// example: true
	IsCorrect bool `json:"is_correct"`

	model.BaseModel // embedded struct with created_at and updated_at
}
