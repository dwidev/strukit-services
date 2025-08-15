package middleware

import (
	"bytes"
	"context"
	"io"
	"strings"
	appContext "strukit-services/pkg/context"
	"strukit-services/pkg/logger"
	"strukit-services/pkg/responses"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		requestId := uuid.New().String()
		c.Header("X-Request-ID", requestId)

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, appContext.RequestIDKey, requestId)
		ctx = context.WithValue(ctx, appContext.IPAddressKey, c.ClientIP())
		ctx = context.WithValue(ctx, appContext.MethodKey, c.Request.Method)
		ctx = context.WithValue(ctx, appContext.PathKey, path)

		projectId := c.Param("project-id")
		if len(projectId) > 0 {
			uuid, err := uuid.Parse(projectId)
			if err != nil {
				c.JSON(400, responses.MessageResponse{StatusCode: 400, Message: err.Error()})
				c.Abort()
				return
			}
			ctx = context.WithValue(ctx, appContext.ProjectID, uuid)
		}

		c.Request = c.Request.WithContext(ctx)

		if strings.Contains(path, "scan/upload") {
			c.Next()
			return
		}

		body, _ := c.GetRawData()
		logger.Log.Request(ctx, string(body)).Info("INCOMING HTTP REQUEST LOG")
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		c.Next()
	}
}
