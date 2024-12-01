package handlers_test

import (
	"api/internal/handlers"
	"api/internal/services"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Мок-сервис для авторизации
func setupAuthService() *services.AuthService {
	return services.NewAuthService(3600, 7200)
}

// Настройка роутеров для тестирования
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	authService := setupAuthService()
	handler := handlers.NewHandler(authService)

	// Роуты для аутентификации
	r.POST("/api/auth/register", handler.RegisterHandler)
	r.POST("/api/auth/login", handler.LoginHandler)
	r.POST("/api/posts", handlers.CreatePostHandler)
	r.GET("/api/posts", handlers.GetPostsHandler)
	r.PUT("/api/posts/:postId", handlers.UpdatePostHandler) // Используем PUT для обновления

	return r
}

func setupRouterWithImageHandlers() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Роуты для работы с изображениями
	r.POST("/api/posts/:postId/images", handlers.AddImageHandler)
	r.DELETE("/api/posts/:postId/images/:imageId", handlers.DeleteImageHandler)

	return r
}

// Тест регистрации пользователя
func TestRegisterHandler(t *testing.T) {
	router := setupRouter()

	body := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"role":     "Author",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Регистрация прошла успешно")
}

// Тест входа пользователя
func TestLoginHandler(t *testing.T) {
	router := setupRouter()

	// Предварительная регистрация пользователя
	authService := setupAuthService()
	authService.RegisterUser("test@example.com", "password123", "Author")

	body := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Вход выполнен успешно")
}

// Тест создания поста
func TestCreatePostHandler(t *testing.T) {
	router := setupRouter()

	body := map[string]string{
		"title":   "Test Post",
		"content": "This is a test post content.",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/posts", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("role", "Author") // Имитируем авторизацию
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Пост успешно создан")
}

// Тест получения всех постов
func TestGetPostsHandler(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/posts", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Список постов")
}

// Тест обновления поста
func TestUpdatePostHandler(t *testing.T) {
	router := setupRouter()

	// Мок создания поста
	postID := "1"
	services.CreatePost(postID, "Test Post")

	body := map[string]string{
		"title":   "Updated Title",
		"content": "Updated Content",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/api/posts/"+postID, bytes.NewBuffer(jsonBody)) // Исправлено на PUT
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("role", "Author") // Имитируем авторизацию
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Пост обновлен")
}

// Тест добавления изображения к посту
func TestAddImageHandler(t *testing.T) {
	router := setupRouterWithImageHandlers()

	// Создаем тестовый пост
	postID := "1"
	services.CreatePost(postID, "Test Post")

	// Данные для запроса
	body := map[string]string{
		"image_url": "https://example.com/image.jpg",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/posts/"+postID+"/images", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("role", "Author") // Имитируем авторизацию
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Изображение добавлено")
}

// Тест ошибки доступа при добавлении изображения (неавторизованный пользователь)
func TestAddImageHandlerUnauthorized(t *testing.T) {
	router := setupRouterWithImageHandlers()

	// Данные для запроса
	body := map[string]string{
		"image_url": "https://example.com/image.jpg",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/api/posts/1/images", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("role", "User") // Неправильная роль
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Доступ запрещен")
}

// Тест удаления изображения из поста
func TestDeleteImageHandler(t *testing.T) {
	router := setupRouterWithImageHandlers()

	// Создаем тестовый пост и добавляем изображение
	postID := "1"
	imageID := "1"
	services.AddImageToPost(postID, imageID)

	req, _ := http.NewRequest("DELETE", "/api/posts/"+postID+"/images/"+imageID, nil)
	req.Header.Set("role", "Admin") // Имитируем авторизацию
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Изображение удалено")
}

// Тест ошибки доступа при удалении изображения (неавторизованный пользователь)
func TestDeleteImageHandlerUnauthorized(t *testing.T) {
	router := setupRouterWithImageHandlers()

	req, _ := http.NewRequest("DELETE", "/api/posts/1/images/1", nil)
	req.Header.Set("role", "User") // Неправильная роль
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Доступ запрещен")
}

// Тест, когда изображение не найдено
func TestDeleteImageHandlerNotFound(t *testing.T) {
	router := setupRouterWithImageHandlers()

	req, _ := http.NewRequest("DELETE", "/api/posts/999/images/999", nil)
	req.Header.Set("role", "Admin") // Имитируем авторизацию
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "не найдены")
}
