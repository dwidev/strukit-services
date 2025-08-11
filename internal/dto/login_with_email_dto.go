package dto

type LoginWithEmailRequest struct {
	Email string `json:"email" validate:"required"`
}
