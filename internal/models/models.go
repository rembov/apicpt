package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Posts = make(map[string]Post)

// Тип данных для постов
type Post struct {
	ID      string
	Title   string
	Content string
	Status  string
}

type User struct {
	Email    string
	Password string
	Role     string
}

type Token struct {
	RefreshToken string
	ExpiresAt    time.Time
}
// JWT ключ
var JwtKey = []byte("super_secret_key")

// Структура для JWT
type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"` 
	jwt.RegisteredClaims
}
var Input struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
var Inputauth struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=Author Reader"`
}
var Inputlogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type AuthService struct {
	users      map[string]User
	tokens     map[string]Token
	tokenTTL   time.Duration
	refreshTTL time.Duration
}