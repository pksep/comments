package model

import (
	"github.com/pksep/location_search_server/internal/modules/shared/model"
	slideModel "github.com/pksep/location_search_server/internal/modules/slides/model"
)

// Project represents an exam project
// @Description Project for exam system
type Project struct {
	// ID of the project
	// example: 123e4567-e89b-12d3-a456-426614174001
	ID string `json:"id"`

	// Name of the project
	// example: Biology Final Exam
	Name string `json:"name"`

	// Description of the project
	// example: Final exam for Biology 101
	Description string `json:"description"`

	// Slides is a list of slides in the project
	Slides []slideModel.Slide `json:"slides,omitempty" gorm:"foreignKey:ProjectID"`

	// ExamDuration is the duration of the exam in minutes
	// example: 60
	ExamDuration int `json:"exam_duration"`

	model.BaseModel // embedded struct with created_at and updated_at
}
