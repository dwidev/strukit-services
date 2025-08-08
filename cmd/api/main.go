package main

import (
	"fmt"
	"strukit-services/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Run(".env.dev")
	router := gin.Default()

	router.Use()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "this server is healthty",
		})
	})

	router.Run(fmt.Sprintf(":%s", config.Env.PORT))
}
