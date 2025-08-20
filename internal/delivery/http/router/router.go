package router

import (
	"strukit-services/internal/delivery/http/middleware"
	"strukit-services/pkg/logger"
	"strukit-services/pkg/token"

	"github.com/gin-gonic/gin"
)

func Run(router *gin.Engine, token *token.Manager, handler *RouterHandler) {
	r := &appRouter{router: router, TokenManager: token, handler: handler}
	r.Build()

}

type appRouter struct {
	router       *gin.Engine
	V1           *gin.RouterGroup
	TokenManager *token.Manager

	handler *RouterHandler
}

func (a *appRouter) Build() {
	if logger.Log == nil {
		panic("please run logger.New() at main.go, before build the router")
	}

	a.router.Use(middleware.SetupContext())
	a.router.Use(middleware.LogRequest())
	a.router.Use(middleware.CatchError())
	a.V1 = a.router.Group("/api/v1")

	a.publicRoutes()
	a.protectedRoutes()
}

func (a *appRouter) publicRoutes() {
	a.router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "this server is healthty",
		})
	})

	auth := a.V1.Group("/auth")
	{
		auth.POST("/login-with-email", a.handler.auth.LoginWithEmail)
		auth.POST("/logout", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "Login dummy",
			})
		})
	}
}

func (a *appRouter) protectedRoutes() {
	a.V1.Use(middleware.Authorization(a.TokenManager))
	user := a.V1.Group("/user")
	{

		user.GET("/profile", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "profile",
			})
		})
		user.PUT("/profile", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "update profile",
			})
		})
	}

	receipt := a.V1.Group("/receipt")
	{
		receipt.POST("/scan/ocr/:project-id", a.handler.receipt.ScanOcr)
		receipt.POST("/scan/upload/:project-id", a.handler.receipt.ScanUpload)
		receipt.GET("/all/:project-id", a.handler.receipt.GetReceiptByProjectID)
		receipt.GET("/detail/:receipt-id", a.handler.receipt.GetDetailReceipt)
		receipt.DELETE("/:receipt-id", a.handler.receipt.OnDelete)
	}

	project := a.V1.Group("/project")
	{
		project.GET("/detail/:project-id", a.handler.project.GetProjectByID)
		project.DELETE("/:id", a.handler.project.SoftDelete)
		project.GET("/all", a.handler.project.All)
		project.POST("/create", a.handler.project.CreateNewProject)
	}

	report := a.V1.Group("/report")
	{

		report.POST("/download/:project-id", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "download  list",
			})
		})
	}
}
