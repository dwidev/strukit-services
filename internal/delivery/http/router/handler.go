package router

import "strukit-services/internal/delivery/http"

func NewHandler(auth *http.AuthHandler) *RouterHandler {
	return &RouterHandler{
		auth: auth,
	}
}

type RouterHandler struct {
	// handler
	auth *http.AuthHandler
}
