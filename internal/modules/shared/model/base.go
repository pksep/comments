package model

import "time"

// BaseModel содержит общие поля для всех сущностей
// @Description Общие поля для всех сущностей: даты создания и обновления
type BaseModel struct {
	// Дата создания сущности
	// example: 2025-09-11T12:34:56Z
	CreatedAt time.Time `json:"created_at"`

	// Дата последнего обновления сущности
	// example: 2025-09-11T12:34:56Z
	UpdatedAt time.Time `json:"updated_at"`
}
