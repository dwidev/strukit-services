package main

import (
	"fmt"
	"strukit-services/internal/delivery/http/router"
	"strukit-services/pkg/config"
	"strukit-services/pkg/db"
	"strukit-services/pkg/logger"

	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println("call init function...")

	config.Run()
	logger.New()

	if config.Env.RuntimeEnv == config.Prod {
		gin.SetMode(gin.ReleaseMode)
	}

	fmt.Println("call init function done")
}

func main() {
	db := db.Open()
	defer db.Close()

	routerEngine := gin.Default()
	router.Run(routerEngine)

	routerEngine.Run(fmt.Sprintf(":%s", config.Env.PORT))
}
