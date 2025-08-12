package middleware

import (
	"errors"
	"strukit-services/pkg/logger"
	"strukit-services/pkg/responses"

	"github.com/gin-gonic/gin"
)

func CatchError() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last()
			logger.Log.Errorf("catched error with : %s", err)

			var appErr *responses.AppError
			if errors.As(err, &appErr) {
				appErr.JSON(ctx)
				return
			}

			responses.ServerError(ctx)
		}

	}
}
