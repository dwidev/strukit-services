package router

import "strukit-services/internal/delivery/http"

func NewHandler(auth *http.AuthHandler, project *http.ProjectHandler, receipt *http.ReceiptHandler, reports *http.ReportHandler) *RouterHandler {
	return &RouterHandler{
		auth:    auth,
		project: project,
		receipt: receipt,
		reports: reports,
	}
}

type RouterHandler struct {
	// handler
	auth    *http.AuthHandler
	project *http.ProjectHandler
	receipt *http.ReceiptHandler
	reports *http.ReportHandler
}
