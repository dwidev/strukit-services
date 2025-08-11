package middleware

import (
	"strukit-services/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return logger.Log.HttpRequestMiddlerware()
}
