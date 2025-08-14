package http

import (
	"io"
	"strukit-services/internal/services"

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

func (r *ReceiptHandler) Scan(c *gin.Context) {
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
