package handlers

import (
	"api/middleware"
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
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if _, exists := users[input.Email]; exists {
		c.JSON(409, gin.H{"error": "Email уже зарегистрирован"})
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

	c.JSON(201, gin.H{"message": "Пользователь успешно зарегистрирован"})
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
		c.JSON(401, gin.H{"error": "Неверные учетные данные"})
		return
	}

	accessToken, err := middleware.GenerateToken(input.Email, tokenTTL)
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
	var input struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for email, token := range tokens {
		if token.RefreshToken == input.RefreshToken && token.ExpiresAt.After(time.Now()) {
			newAccessToken, err := middleware.GenerateToken(email, tokenTTL)
			if err != nil {
				c.JSON(500, gin.H{"error": "Ошибка создания токена"})
				return
			}
			c.JSON(200, gin.H{"accessToken": newAccessToken})
			return
		}
	}

	c.JSON(401, gin.H{"error": "Неверный или истекший токен"})
}
