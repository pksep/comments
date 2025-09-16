package model

import (
	"github.com/pksep/location_search_server/internal/modules/shared/model"
)

// User представляет пользователя системы
// @Description Пользователь экзаменационной системы
type User struct {
	// ID пользователя
	// example: 123e4567-e89b-12d3-a456-426614174000
	ID string `json:"id"`

	// Инициалы пользователя
	// example: AB
	Initials string `json:"initials"`

	model.BaseModel // встроенная структура
}
