package app

import (
	"strukit-services/internal/delivery/http"
	"strukit-services/internal/delivery/http/router"
	"strukit-services/internal/services"
	"strukit-services/pkg/config"
	"strukit-services/pkg/db"
	"strukit-services/pkg/token"

	"github.com/gin-gonic/gin"
)

type BootstrapConfig struct {
	RouterEngine *gin.Engine
}

func Bootstrap(cfg *BootstrapConfig) {
	db := db.Open()
	defer db.Close()
	token := token.Generator(config.Env.JWT_ACCESS_SECRET, config.Env.JWT_REFRESH_SECRET)

	// SERVICE
	authService := services.NewAuth(token)

	// HANDLER
	authHandler := http.NewAuth(authService)

	// ROUTE HANDLER
	routerHandler := router.NewHandler(authHandler)

	// run router
	router.Run(cfg.RouterEngine, routerHandler)

}
