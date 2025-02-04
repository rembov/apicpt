package models

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Модель токена
type Token struct {
	gorm.Model
	UserID       uint
	RefreshToken string
	ExpiresAt    time.Time
	User         User `gorm:"foreignKey:UserID"`
}

// JWT ключ
var JwtKey = []byte("super_secret_key")

// Структура для JWT
type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// Структуры для входных данных
type Input struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type InputAuth struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=Author Reader"`
}

type InputLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Сервис аутентификации
type AuthService struct {
	db         *gorm.DB
	tokenTTL   time.Duration
	refreshTTL time.Duration
}

// Структура подключения к БД
type Database struct {
	DB *gorm.DB
}

// Инициализация базы данных
func InitDB() *gorm.DB {
	dsn := "host=localhost user=db password=db1234 dbname=postgres port=5455 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Автоматическая миграция моделей
	err = db.AutoMigrate(&Post{}, &User{}, &Token{})
	if err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}

	return db
}

// Создание нового сервиса аутентификации
func NewAuthService(db *gorm.DB, tokenTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		db:         db,
		tokenTTL:   tokenTTL,
		refreshTTL: refreshTTL,
	}
}
