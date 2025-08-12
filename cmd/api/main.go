package main

import (
	"fmt"
	"strukit-services/internal/app"
	"strukit-services/pkg/config"
	"strukit-services/pkg/logger"

	"github.com/gin-gonic/gin"
)

func init() {
	config.Run()
	logger.New()

	if config.Env.RuntimeEnv == config.Prod {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	routerEngine := gin.Default()
	cfg := &app.BootstrapConfig{
		RouterEngine: routerEngine,
	}

	close := app.Bootstrap(cfg)
	defer close()

	if err := routerEngine.Run(fmt.Sprintf(":%s", config.Env.PORT)); err != nil {
		logger.Log.Fatalf("Error when starting server : %s", err)
	}
}
