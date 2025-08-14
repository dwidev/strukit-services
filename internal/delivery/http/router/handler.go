package router

import "strukit-services/internal/delivery/http"

func NewHandler(auth *http.AuthHandler, project *http.ProjectHandler, receipt *http.ReceiptHandler) *RouterHandler {
	return &RouterHandler{
		auth:    auth,
		project: project,
		receipt: receipt,
	}
}

type RouterHandler struct {
	// handler
	auth    *http.AuthHandler
	project *http.ProjectHandler
	receipt *http.ReceiptHandler
}
