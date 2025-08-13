package app

import (
	"strukit-services/internal/delivery/http"
	"strukit-services/internal/delivery/http/router"
	"strukit-services/internal/repository"
	"strukit-services/internal/services"
	"strukit-services/pkg/config"
	"strukit-services/pkg/db"
	"strukit-services/pkg/token"
	"strukit-services/pkg/validator"

	"github.com/gin-gonic/gin"
)

type BootstrapConfig struct {
	RouterEngine *gin.Engine
}

func Bootstrap(cfg *BootstrapConfig) func() {
	db := db.Open()

	token := token.NewManager(config.Env.JWT_ACCESS_SECRET, config.Env.JWT_REFRESH_SECRET)
	appValidator := validator.Run()

	// BASE
	baseRepo := repository.NewBase(db.Instance())
	baseHandler := http.NewBase(appValidator)

	// REPOSITORY
	userRepo := repository.NewUser(baseRepo)
	projectRepo := repository.NewProject(baseRepo)

	// SERVICE
	authService := services.NewAuth(token, userRepo)
	projectService := services.NewProject(projectRepo)

	// HANDLER
	authHandler := http.NewAuth(authService)
	projectHandler := http.NewProject(baseHandler, projectService)

	// ROUTE HANDLER
	routerHandler := router.NewHandler(authHandler, projectHandler)

	// run router
	router.Run(cfg.RouterEngine, token, routerHandler)

	close := func() {
		db.Close()
	}
	return close
}
