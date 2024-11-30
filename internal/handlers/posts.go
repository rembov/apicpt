package handlers

import (
	"api/internal/models"
	"api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler создает новый пост
func CreatePostHandler(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "Author" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Пользователь не является автором"})
		return
	}

	if err := c.ShouldBindJSON(&models.Input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postID, err := services.CreatePost(models.Input.Title, models.Input.Content)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Пост успешно создан", "id": postID})
}

// GetPostsHandler возвращает список постов
func GetPostsHandler(c *gin.Context) {
	posts := services.GetPublishedPosts()
	c.JSON(http.StatusOK, posts)
}

// UpdatePostHandler обновляет пост
func UpdatePostHandler(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "Author" && role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	postID := c.Param("id")
	 
	if err := c.ShouldBindJSON(&models.Input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdatePost(postID, models.Input.Title, models.Input.Content); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пост обновлен"})
}

// PublishPostHandler публикует пост
func PublishPostHandler(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "Admin" && role != "Author" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещён"})
		return
	}

	postID := c.Param("id")
	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.PublishPost(postID, input.Status); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пост опубликован"})
}
