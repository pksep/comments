package dto

type CreateUserDTO struct {
	Initials string `json:"initials" binding:"required"`
}