package main

import (
	"fmt"
	"strukit-services/internal/app"
	"strukit-services/pkg/config"
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
	routerEngine := gin.Default()
	cfg := &app.BootstrapConfig{
		RouterEngine: routerEngine,
	}
	app.Bootstrap(cfg)
	routerEngine.Run(fmt.Sprintf(":%s", config.Env.PORT))
}
