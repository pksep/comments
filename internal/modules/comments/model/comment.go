package model

import (
    shared "github.com/pksep/comments/internal/modules/shared/model"
)

// Comment is a reusable comment entity that can be attached to any domain entity
// by specifying entity type and entity id.
type Comment struct {
    ID          string       `json:"id"`
    EntityType  string       `json:"entity_type"`
    EntityID    string       `json:"entity_id"`
    AuthorID    string       `json:"author_id"`
    Content     string       `json:"content"`
    ParentID    *string      `json:"parent_id,omitempty"`
    shared.BaseModel
}


