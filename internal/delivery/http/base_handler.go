package http

import "strukit-services/pkg/validator"

func NewBase(validator *validator.AppValidator) *BaseHandler {
	return &BaseHandler{
		AppValidator: validator,
	}
}

type BaseHandler struct {
	*validator.AppValidator
}
