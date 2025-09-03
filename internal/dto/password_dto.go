package dto

type CreatePasswordDTO struct {
	Password string `json:"password" validate:"required"`
}
