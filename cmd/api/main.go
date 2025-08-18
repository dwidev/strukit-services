package main

import (
	"fmt"
	"os"
	"strukit-services/internal/app"
	"strukit-services/pkg/config"
	"strukit-services/pkg/constant"
	"strukit-services/pkg/logger"

	"github.com/gin-gonic/gin"
)

func init() {
	env := os.Getenv("GO_ENV")
	logger.New(&logger.Config{Env: constant.Environment(env)})
	config.Run(env)

	if config.Env.RuntimeEnv == constant.Prod {
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
