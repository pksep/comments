package dto

type UpdateUserDTO struct {
	Initials string `json:"initials" binding:"required"`
}