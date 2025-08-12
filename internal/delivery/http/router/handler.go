package router

import "strukit-services/internal/delivery/http"

func NewHandler(auth *http.AuthHandler, project *http.ProjectHandler) *RouterHandler {
	return &RouterHandler{
		auth:    auth,
		project: project,
	}
}

type RouterHandler struct {
	// handler
	auth    *http.AuthHandler
	project *http.ProjectHandler
}
