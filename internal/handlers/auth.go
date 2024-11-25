package handlers

import (
	"api/internal/middleware"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат email"})
		return
	}

	if _, exists := users[input.Email]; exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Email уже зарегистрирован"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка"})
		return
	}

	users[input.Email] = User{
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	c.JSON(http.StatusOK, gin.H{"message": "Регистрация прошла успешно"})
}

func LoginHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := users[input.Email]
	if !exists || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	accessToken, err := middleware.GenerateToken(input.Email, "User", tokenTTL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания токена"})
		return
	}

	refreshToken := uuid.NewString()
	tokens[input.Email] = Token{
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(refreshTTL),
	}

	// Установка access токена в куки
	c.SetCookie("accessToken", accessToken, int(tokenTTL.Seconds()), "/", "", false, true)

	// Установка refresh токена в куки
	c.SetCookie("refreshToken", refreshToken, int(refreshTTL.Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Вход выполнен успешно"})
}

func RefreshTokenHandler(c *gin.Context) {
	// Получение refresh токена из куки
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh токен не найден"})
		return
	}

	for email, token := range tokens {
		if token.RefreshToken == refreshToken && token.ExpiresAt.After(time.Now()) {
			newAccessToken, err := middleware.GenerateToken(email, "User", tokenTTL)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":       "Ошибка генерации токена",
					"description": err.Error(),
				})
				return
			}

			// Обновление access токена в куки
			c.SetCookie("accessToken", newAccessToken, int(tokenTTL.Seconds()), "/", "", false, true)

			c.JSON(http.StatusOK, gin.H{"message": "Токен успешно обновлён"})
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh токен недействителен"})
}
