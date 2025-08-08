package router

import "github.com/gin-gonic/gin"

func Run(router *gin.Engine) {
	r := &appRouter{router: router}
	r.BuidlV1().Build()
}

type appRouter struct {
	router *gin.Engine
	V1     *gin.RouterGroup

	// controller
}

func (a *appRouter) BuidlV1() *appRouter {
	a.V1 = a.router.Group("/api/v1")
	return a
}

func (a *appRouter) Build() {
	a.router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "this server is healthty",
		})
	})

	auth := a.V1.Group("/auth")
	{

		auth.POST("/login", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "Login dummy",
			})
		})
		auth.POST("/logout", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "Login dummy",
			})
		})
	}

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

		project.GET("/create", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "receipt list",
			})
		})
		project.GET("/delete", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "receipt list",
			})
		})
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
