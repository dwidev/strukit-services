package middleware

import (
	"context"
	"net/http"
	"strings"
	"strukit-services/pkg/config"
	appContext "strukit-services/pkg/context"
	"strukit-services/pkg/token"

	"github.com/gin-gonic/gin"
)

func Authorization(token *token.Token) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqCtx := ctx.Request.Context()
		accessTokenReq := ctx.GetHeader("access-token")

		if len(accessTokenReq) <= 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "access-token header required"})
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(accessTokenReq, " ")
		if len(tokenParts) < 2 || tokenParts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format. Use 'Bearer <token>'"})
			ctx.Abort()
			return
		}

		accessToken := tokenParts[1]
		tokenParse, err := token.Parse(accessToken, config.Env.JWT_ACCESS_SECRET)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		newCtx := context.WithValue(reqCtx, appContext.UserIDKey, tokenParse.UserID)
		ctx.Request = ctx.Request.WithContext(newCtx)
		ctx.Next()
	}
}
