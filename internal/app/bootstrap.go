package app

import (
	"strukit-services/internal/delivery/http"
	"strukit-services/internal/delivery/http/router"
	"strukit-services/internal/repository"
	"strukit-services/internal/services"
	"strukit-services/pkg/config"
	"strukit-services/pkg/db"
	"strukit-services/pkg/token"

	"github.com/gin-gonic/gin"
)

type BootstrapConfig struct {
	RouterEngine *gin.Engine
}

func Bootstrap(cfg *BootstrapConfig) func() {
	db := db.Open()

	token := token.Generator(config.Env.JWT_ACCESS_SECRET, config.Env.JWT_REFRESH_SECRET)

	// BASE
	baseRepo := repository.NewBase(db.Instance())

	// REPOSITORY
	userRepo := repository.NewUser(baseRepo)

	// SERVICE
	authService := services.NewAuth(token, userRepo)

	// HANDLER
	authHandler := http.NewAuth(authService)

	// ROUTE HANDLER
	routerHandler := router.NewHandler(authHandler)

	// run router
	router.Run(cfg.RouterEngine, token, routerHandler)

	close := func() {
		db.Close()
	}
	return close
}
