package middleware

import "github.com/gin-gonic/gin"

func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
