package model

import (
	"github.com/pksep/location_search_server/internal/modules/shared/model"
)

// Slide represents a slide in an exam project
// @Description Slide in an exam project
type Slide struct {
	// ID of the slide
	// example: 123e4567-e89b-12d3-a456-426614174002
	ID string `json:"id"`

	// ProjectID is the ID of the project this slide belongs to
	// example: 123e4567-e89b-12d3-a456-426614174001
	ProjectID string `json:"project_id" gorm:"index"`

	// Path to the slide image
	// example: /images/slides/slide1.jpg
	Path string `json:"path"`

	// Text content of the slide
	// example: This is a sample question slide
	Text string `json:"text"`

	model.BaseModel // embedded struct with created_at and updated_at
}
