package router

import (
	"strukit-services/internal/delivery/http/middleware"
	"strukit-services/pkg/logger"
	"strukit-services/pkg/token"

	"github.com/gin-gonic/gin"
)

func Run(router *gin.Engine, token *token.Token, handler *RouterHandler) {
	r := &appRouter{router: router, Token: token, handler: handler}
	r.Build()

}

type appRouter struct {
	router *gin.Engine
	V1     *gin.RouterGroup
	Token  *token.Token

	handler *RouterHandler
}

func (a *appRouter) Build() {
	if logger.Log == nil {
		panic("please run logger.New() at main.go")
	}

	a.router.Use(middleware.Logging())
	a.router.Use(middleware.CatchError())
	a.V1 = a.router.Group("/api/v1")

	a.PublicRoute()
	a.AuthRoute()
}

func (a *appRouter) PublicRoute() {
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

func (a *appRouter) AuthRoute() {
	a.V1.Use(middleware.Authorization(a.Token))
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
		receipt.DELETE("/:id", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "receipt detail",
			})
		})
		receipt.POST("/scan", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "receipt scan",
			})
		})
		receipt.GET("/list", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "receipt list",
			})
		})
		receipt.GET("/detail/:id", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "receipt detail",
			})
		})

	}

	project := a.V1.Group("/project")
	{
		project.DELETE("/:id", a.handler.project.SoftDelete)
		project.GET("/all", a.handler.project.SoftDelete)
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
