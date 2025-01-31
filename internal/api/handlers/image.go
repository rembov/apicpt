package handlers

import (
	"apicpt/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	role, exists := c.Get("role")
	if !exists || (role != "Admin" && role != "Author") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	postID := c.Param("postId")
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось загрузить файл"})
		return
	}

	// Сохраняем файл
	path := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения файла"})
		return
	}

	// Получаем объект базы данных из контекста
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}

	// Добавляем изображение к посту
	if err := services.AddImageToPost(db.(*gorm.DB), postID, path); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Изображение добавлено", "path": path})
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
	role, exists := c.Get("role")
	if !exists || role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен"})
		return
	}

	postID := c.Param("postId")
	imageID := c.Param("imageId")

	// Получаем объект базы данных из контекста
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}

	// Удаляем изображение
	if err := services.RemoveImageFromPost(db.(*gorm.DB), postID, imageID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Изображение удалено"})
}
