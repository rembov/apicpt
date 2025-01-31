package api

import (
	"apicpt/internal/api/handlers"
	"apicpt/internal/middleware"
	"apicpt/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, authService *services.AuthService) *gin.Engine {
	r := gin.Default()

	authHandler := handlers.NewHandler(authService)
	r.Use(func(c *gin.Context) {
		c.Set("db", db) // Передаем базу данных через контекст
		c.Next()
	})

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.RegisterHandler)
		auth.POST("/login", authHandler.LoginHandler)
		auth.POST("/refresh-token", authHandler.RefreshTokenHandler)
	}

	// Эндпоинты управления постами
	postsGroup := r.Group("/api/posts")
	postsGroup.Use(middleware.AuthMiddleware)
	{
		postsGroup.POST("", handlers.CreatePostHandler)
		postsGroup.GET("", handlers.GetPostsHandler)
		postsGroup.PUT("/:id", handlers.UpdatePostHandler)
		postsGroup.PATCH("/:id/status", handlers.PublishPostHandler)

		// Группа для работы с изображениями
		images := postsGroup.Group("/:id/images")
		{
			images.POST("", handlers.AddImageHandler)
			images.DELETE("/:imageId", handlers.DeleteImageHandler)
		}
	}

	return r
}
