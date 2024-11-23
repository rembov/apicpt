package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddImageHandler добавляет изображение к посту
func AddImageHandler(c *gin.Context) {
	// Логика загрузки изображения
	c.JSON(http.StatusCreated, gin.H{"message": "Изображение добавлено"})
}

// DeleteImageHandler удаляет изображение из поста
func DeleteImageHandler(c *gin.Context) {
	// Логика удаления изображения
	c.JSON(http.StatusOK, gin.H{"message": "Изображение удалено"})
}
