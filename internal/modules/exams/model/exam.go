package model

import (
	answersModel "github.com/pksep/location_search_server/internal/modules/answers/model"
	sharemodel "github.com/pksep/location_search_server/internal/modules/shared/model"
)

// Exam represents an exam taken by a user
// @Description Exam taken by a user
type Exam struct {
	// ID of the exam
	// example: 123e4567-e89b-12d3-a456-426614174003
	ID string `json:"id"`

	// ProjectID is the ID of the project this exam is based on
	// example: 123e4567-e89b-12d3-a456-426614174001
	ProjectID string `json:"project_id" gorm:"index"`

	// UserID is the ID of the user taking the exam
	// example: 123e4567-e89b-12d3-a456-426614174000
	UserID string `json:"user_id" gorm:"index"`

	// IsActive indicates if the exam is currently in progress
	// example: true
	IsActive bool `json:"is_active" gorm:"default:true"`

	// Answers is a list of answers provided during the exam
	Answers []answersModel.Answer `json:"answers,omitempty" gorm:"foreignKey:ExamID"`

	sharemodel.BaseModel // embedded struct with created_at and updated_at
}
