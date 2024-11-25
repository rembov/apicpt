package router

import (
	"api/internal/handlers"
	"api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter настраивает маршруты
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Эндпоинты аутентификации
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.RegisterHandler)
		auth.POST("/login", handlers.LoginHandler)
		auth.POST("/refresh-token", handlers.RefreshTokenHandler)
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
