package handlers

import (
	"api/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Имитация базы данных
var (
	users      = make(map[string]User)
	tokens     = make(map[string]Token)
	tokenTTL   = time.Hour * 2
	refreshTTL = time.Hour * 24 * 7
)

type User struct {
	Email    string
	Password string
	Role     string
}

type Token struct {
	RefreshToken string
	ExpiresAt    time.Time
}

func RegisterHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Role     string `json:"role" binding:"required,oneof=Author Reader"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Неверный формат email"})
		return
	}

	if _, exists := users[input.Email]; exists {
		c.JSON(403, gin.H{"error": "Email уже зарегистрирован"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка хэширования пароля"})
		return
	}

	users[input.Email] = User{
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	c.JSON(200, gin.H{"message": "Регистрация прошла успешно"})
}

func LoginHandler(c *gin.Context) {

	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, exists := users[input.Email]
	if !exists || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		c.JSON(403, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	accessToken, err := middleware.GenerateToken(input.Email, "User", tokenTTL)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка создания токена"})
		return
	}

	refreshToken := uuid.NewString()
	tokens[input.Email] = Token{
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(refreshTTL),
	}

	c.JSON(200, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func RefreshTokenHandler(c *gin.Context) {
	// Входящие данные (структура, которую ожидает документация)
	var input struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	// Проверяем корректность входящих данных
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Токен недействителен",
			"description": err.Error(),
		})
		return
	}

	for email, token := range tokens {
		if token.RefreshToken == input.RefreshToken && token.ExpiresAt.After(time.Now()) {
			// Если токен найден и не истёк, создаём новый accessToken
			newAccessToken, err := middleware.GenerateToken(email, "User", time.Hour*2)
			if err != nil {
				// Обработка ошибки генерации токена
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":       "Ошибка генерации токена",
					"description": err.Error(),
				})
				return
			}

			// Успешный респонс с новым токеном
			c.JSON(http.StatusOK, gin.H{
				"accessToken": newAccessToken,
			})
			return
		}
	}

	// Если токен не найден или истёк
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Токен обновления недействителен",
	})
}
