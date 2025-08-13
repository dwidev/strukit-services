package middleware

import (
	"context"
	appContext "strukit-services/pkg/context"
	"strukit-services/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := logger.Log
		requestId := uuid.New().String()
		c.Header("X-Request-ID", requestId)

		ctx := context.WithValue(c.Request.Context(), appContext.RequestIDKey, requestId)
		c.Request = c.Request.WithContext(ctx)

		data := l.LogRequest(c, requestId)
		l.WithFields(data).Info("HTTP REQUEST LOG")
		c.Next()
	}
}
