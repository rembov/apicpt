package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Имитация базы данных для постов
var posts = make(map[string]Post)

// Тип данных для постов
type Post struct {
	ID      string
	Title   string
	Content string
	Status  string
}

// CreatePostHandler создает новый пост
func CreatePostHandler(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "Author" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Пользователь не является автором",
		})
		return
	}
	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	
	for _, post := range posts {
        if post.Title == input.Title {
            c.JSON(http.StatusConflict, gin.H{ 
                "error": "Уникальный ключ уже исползьзовался",
            })
            return
        }
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postID := uuid.NewString()
	posts[postID] = Post{
		ID:      postID,
		Title:   input.Title,
		Content: input.Content,
		Status:  "Draft",
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Пост успешно создан", "id": postID})
}

// GetPostsHandler возвращает список постов
func GetPostsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, posts)
}

// UpdatePostHandler обновляет пост
func UpdatePostHandler(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "Author" && role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Доступ запрещен",
		})
		return
	}
	postID := c.Param("id")
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, exists := posts[postID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пост не найден"})
		return
	}

	post.Title = input.Title
	post.Content = input.Content
	posts[postID] = post

	c.JSON(http.StatusOK, gin.H{"message": "Пост обновлен"})
}

// PublishPostHandler публикует пост
func PublishPostHandler(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "Admin" && role != "Author" {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "Доступ запрещён: публикация поста доступна только автору или администратору.",
        })
        return
    }
	postID := c.Param("id")
	var input struct {
		Status string `json:"status" binding:"required,oneof=Published"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, exists := posts[postID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пост не найден"})
		return
	}

	post.Status = input.Status
	posts[postID] = post

	c.JSON(http.StatusOK, gin.H{"message": "Пост успешно создан"})
}
