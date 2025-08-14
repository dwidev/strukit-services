package http

import (
	"io"
	"net/http"
	"strukit-services/internal/services"
	"strukit-services/pkg/logger"

	"github.com/gin-gonic/gin"
)

func NewReceipt(ReceiptService *services.ReceiptService) *ReceiptHandler {
	return &ReceiptHandler{
		ReceiptService: ReceiptService,
	}
}

type ReceiptHandler struct {
	BaseHandler
	*services.ReceiptService
}

func (r *ReceiptHandler) ScanOcr(c *gin.Context) {
	ctx := c.Request.Context()

	type Request struct {
		Raw string `json:"rawOcr" validate:"required"`
	}

	request := Request{}
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Log.Errorf("error binding request %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := r.ReceiptService.ScanFromOCR(ctx, request.Raw)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, res)
}

func (r *ReceiptHandler) ScanUpload(c *gin.Context) {
	ctx := c.Request.Context()
	image, err := c.FormFile("image")

	if err != nil {
		c.Error(err)
		return
	}

	openedFile, err := image.Open()
	if err != nil {
		c.Error(err)
		return
	}
	defer openedFile.Close()

	imageData, err := io.ReadAll(openedFile)
	if err != nil {
		c.Error(err)
		return
	}

	res, err := r.ReceiptService.Scan(ctx, imageData)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, res)
}
