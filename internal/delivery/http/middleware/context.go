package middleware

import (
	"context"
	appContext "strukit-services/pkg/context"
	"strukit-services/pkg/responses"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetupContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		requestId := uuid.New().String()
		c.Header("X-Request-ID", requestId)

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, appContext.RequestIDKey, requestId)
		ctx = context.WithValue(ctx, appContext.IPAddressKey, c.ClientIP())
		ctx = context.WithValue(ctx, appContext.MethodKey, c.Request.Method)
		ctx = context.WithValue(ctx, appContext.PathKey, path)

		if projectId := c.Param("project-id"); projectId != "" {
			uuid, err := uuid.Parse(projectId)
			if err != nil {
				c.JSON(400, responses.MessageResponse{StatusCode: 400, Message: err.Error()})
				c.Abort()
				return
			}
			ctx = context.WithValue(ctx, appContext.ProjectID, uuid)
		}

		if receiptID := c.Param("receipt-id"); receiptID != "" {
			uuid, err := uuid.Parse(receiptID)
			if err != nil {
				c.JSON(400, responses.MessageResponse{StatusCode: 400, Message: err.Error()})
				c.Abort()
				return
			}
			ctx = context.WithValue(ctx, appContext.ReceiptIDKey, uuid)
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
