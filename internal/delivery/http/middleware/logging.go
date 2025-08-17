package middleware

import (
	"bytes"
	"io"
	"strings"
	"strukit-services/pkg/logger"

	"github.com/gin-gonic/gin"
)

func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		path := c.Request.URL.Path

		if strings.Contains(path, "scan/upload") {
			logger.Log.Request(ctx, "[FILE UPLOAD]").Info("INCOMING HTTP REQUEST LOG")
			c.Next()
			return
		}

		if body, err := c.GetRawData(); err == nil {
			logger.Log.Request(ctx, string(body)).Info("INCOMING HTTP REQUEST LOG")
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		c.Next()
	}
}
