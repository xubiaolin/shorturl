package main

import (
	"log"

	"shorturl/handlers"
	"shorturl/middleware"
	"shorturl/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	if err := models.InitDB("shorturl.db"); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化默认用户
	if err := handlers.InitDefaultUser(); err != nil {
		log.Fatalf("Failed to initialize default user: %v", err)
	}

	// 创建 Gin 引擎
	r := gin.Default()

	// 初始化 handlers
	shortURLHandler := handlers.NewShortURLHandler()
	redirectHandler := handlers.NewRedirectHandler()
	authHandler := handlers.NewAuthHandler()

	// 设置路由
	// 健康检查接口（无需鉴权）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 登录接口（无需鉴权）
	r.POST("/api/v1/login", authHandler.Login)

	// 短链跳转接口（无需鉴权）
	r.GET("/:code", redirectHandler.Redirect)

	// API v1 路由组（需要鉴权）
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		// 认证相关接口
		api.POST("/change-password", authHandler.ChangePassword)

		// 短链管理接口
		api.POST("/shorturls", shortURLHandler.CreateShortURL)
		api.PUT("/shorturls/:id", shortURLHandler.UpdateShortURL)
		api.DELETE("/shorturls/:id", shortURLHandler.DeleteShortURL)
		api.GET("/shorturls/:id", shortURLHandler.GetShortURL)
		api.GET("/shorturls", shortURLHandler.ListShortURLs)
		api.GET("/shorturls/:id/stats", shortURLHandler.StatsShortURL)
	}

	// 静态文件服务
	r.Static("/static", "./static")

	// 前端页面路由
	r.GET("/login", func(c *gin.Context) {
		c.File("./static/login.html")
	})
	r.GET("/change-password", func(c *gin.Context) {
		c.File("./static/change-password.html")
	})
	r.GET("/dashboard", func(c *gin.Context) {
		c.File("./static/dashboard.html")
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/login")
	})

	// 启动服务
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
