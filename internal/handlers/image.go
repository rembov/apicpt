package handlers

import (
	"api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddImageHandler добавляет изображение к посту
func AddImageHandler(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "Admin" && role != "Author" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	postID := c.Param("postId")
	var input struct {
		ImageURL string `json:"image_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.AddImageToPost(postID, input.ImageURL); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Изображение добавлено"})
}

// DeleteImageHandler удаляет изображение из поста
func DeleteImageHandler(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	postID := c.Param("postId")

	if err := services.RemoveImageFromPost(postID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Изображение удалено"})
}
