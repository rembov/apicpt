package handlers

import (
	"api/internal/models"
	"api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler
// @Summary Создание поста
// @Tags Управление постами
// @Description Создает новый пост с заголовком и содержимым
// @Accept multipart/form-data
// @Produce json
// @Param idempotencyKey formData string true "Уникальный ключ"
// @Param title formData string true "Заголовок поста"
// @Param content formData string true "Содержимое поста"
// @Security bearerAuth
// @Success 201 {string} string "Пост успешно создан"
// @Failure 403 {string} string "Пользователь не является автором"
// @Failure 409 {string} string "Уникальный ключ уже использовался"
// @Router /api/posts [post]
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

// GetPostsHandler
// @Summary Получение списка постов
// @Tags Получение постов
// @Description Возвращает список постов
// @Produce json
// @Security bearerAuth
// @Success 200 {array} string "Список постов"
// @Router /api/posts [get]
func GetPostsHandler(c *gin.Context) {
	posts := services.GetPublishedPosts()
	c.JSON(http.StatusOK, posts)
}

// UpdatePostHandler
// @Summary Редактирование поста
// @Tags Управление постами
// @Description Редактирует пост по указанному идентификатору
// @Accept multipart/form-data
// @Produce json
// @Param postId path int true "Идентификатор поста"
// @Param title formData string false "Новый заголовок поста"
// @Param content formData string false "Новое содержимое поста"
// @Security bearerAuth
// @Success 200 {string} string "Пост успешно обновлен"
// @Failure 404 {string} string "Пост не найден"
// @Failure 403 {string} string "Доступ запрещен"
// @Router /api/posts/{postId} [post]
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

// PublishPostHandler
// @Summary Публикация поста
// @Tags Управление постами
// @Description Публикует пост по указанному идентификатору
// @Accept multipart/form-data
// @Produce json
// @Param postId path int true "Идентификатор поста"
// @Param status formData string true "Статус поста (Published)"
// @Security bearerAuth
// @Success 200 {string} string "Пост успешно опубликован"
// @Failure 400 {string} string "Неверное значение статуса"
// @Failure 404 {string} string "Пост не найден"
// @Router /api/posts/{postId}/status [patch]
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
