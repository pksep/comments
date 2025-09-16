package dto

type CreateCommentDTO struct {
    EntityType string  `json:"entity_type"`
    EntityID   string  `json:"entity_id"`
    AuthorID   string  `json:"author_id"`
    Content    string  `json:"content"`
    ParentID   *string `json:"parent_id,omitempty"`
}


