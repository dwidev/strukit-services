package app

import (
	"context"
	"fmt"
	"strukit-services/internal/delivery/http"
	"strukit-services/internal/delivery/http/router"
	"strukit-services/internal/repository"
	"strukit-services/internal/services"
	"strukit-services/pkg/config"
	"strukit-services/pkg/db"
	"strukit-services/pkg/llm"
	"strukit-services/pkg/token"
	"strukit-services/pkg/validator"

	"github.com/gin-gonic/gin"
)

type BootstrapConfig struct {
	RouterEngine *gin.Engine
}

func Bootstrap(cfg *BootstrapConfig) func() {
	defaultCtx := context.Background()
	db := db.Open()

	token := token.NewManager(config.Env.JWT_ACCESS_SECRET, config.Env.JWT_REFRESH_SECRET)
	appValidator := validator.Run()
	llm, err := llm.Run(defaultCtx)
	if err != nil {
		panic(fmt.Sprintf("panic when run llm Manager with error : %s", err))
	}

	// BASE
	baseRepo := repository.NewBase(db.Instance())
	baseHandler := http.NewBase(appValidator)

	// REPOSITORY
	userRepo := repository.NewUser(baseRepo)
	projectRepo := repository.NewProject(baseRepo)

	// SERVICE
	authService := services.NewAuth(token, userRepo)
	projectService := services.NewProject(projectRepo)
	receiptService := services.NewReceipt(llm)

	// HANDLER
	authHandler := http.NewAuth(authService)
	projectHandler := http.NewProject(baseHandler, projectService)
	receiptHandler := http.NewReceipt(receiptService)

	// ROUTE HANDLER
	routerHandler := router.NewHandler(authHandler, projectHandler, receiptHandler)

	// run router
	router.Run(cfg.RouterEngine, token, routerHandler)

	close := func() {
		db.Close()
	}
	return close
}
