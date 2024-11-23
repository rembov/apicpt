package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT ключ
var jwtKey = []byte("super_secret_key")

// Структура для JWT
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// AuthMiddleware проверяет авторизацию
func AuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется токен авторизации"})
		c.Abort()
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
		c.Abort()
		return
	}

	c.Set("email", claims.Email)
	c.Next()
}

// GenerateToken генерирует новый токен
func GenerateToken(email string, duration time.Duration) (string, error) {
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
