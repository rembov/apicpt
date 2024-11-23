package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddImageHandler добавляет изображение к посту
func AddImageHandler(c *gin.Context) {
	postID := c.Param("postId")

	// Проверяем, существует ли пост
	post, exists := posts[postID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пост или картинка не найдены"})
		return
	}

	post.Content += "[Image:]"
	posts[postID] = post

	c.JSON(http.StatusCreated, gin.H{"message": "Изображение добавлено"})
}

// DeleteImageHandler удаляет изображение из поста
func DeleteImageHandler(c *gin.Context) {
	postID := c.Param("postId")

	// Проверяем, существует ли пост
	post, exists := posts[postID]
	if !exists {
		c.JSON(404, gin.H{"error": "Пост не найден"})
		return
	}

	// Удаляем привязку к изображению из поста

	post.Content = "[Image removed]"
	posts[postID] = post

	c.JSON(http.StatusOK, gin.H{"message": "Изображение удалено"})
}
