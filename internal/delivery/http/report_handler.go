package http

import (
	"fmt"
	"strukit-services/internal/services"

	"github.com/gin-gonic/gin"
)

func NewReport(base *BaseHandler, reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{
		BaseHandler:   base,
		reportService: reportService,
	}
}

type ReportHandler struct {
	*BaseHandler
	reportService *services.ReportService
}

func (a ReportHandler) DownloadExcelFile(c *gin.Context) {

	file, err := a.reportService.DownloadExcelFile(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", "report.xls"))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	_ = file.Write(c.Writer)

}
