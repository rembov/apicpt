package handlers

import (
	"api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)
// AddImageHandler
// @Summary Добавление картинки к посту
// @Tags Управление постами
// @Description Добавляет изображение к посту
// @Accept multipart/form-data
// @Produce json
// @Param postId path int true "Идентификатор поста"
// @Param image formData file true "Изображение"
// @Security bearerAuth
// @Success 201 {string} string "Картинка добавлена к посту"
// @Failure 404 {string} string "Пост не найден"
// @Failure 403 {string} string "Доступ запрещен"
// @Router /api/posts/{postId}/images [post]
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
// DeleteImageHandler
// @Summary Удаление картинки из поста
// @Tags Управление постами
// @Description Удаляет изображение из поста по идентификатору
// @Param postId path int true "Идентификатор поста"
// @Param imageId path int true "Идентификатор картинки"
// @Security bearerAuth
// @Success 200 {string} string "Картинка успешно удалена"
// @Failure 404 {string} string "Пост или картинка не найдены"
// @Failure 403 {string} string "Доступ запрещён"
// @Router /api/posts/{postId}/images/{imageId} [delete]
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
