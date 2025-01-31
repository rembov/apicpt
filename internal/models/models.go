package models

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Модель поста
type Post struct {
	gorm.Model
	ID        string `gorm:"primaryKey"`
	Title     string
	Content   string
	Status    string
	AuthorID  uint
	Author    User `gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Модель пользователя
type User struct {
	gorm.Model
	Email     string `gorm:"uniqueIndex"`
	Password  string
	Role      string
	Posts     []Post  `gorm:"foreignKey:AuthorID"`
	Tokens    []Token `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

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

// Методы для работы с постами
func (db *Database) CreatePost(post *Post) error {
	return db.DB.Create(post).Error
}

func (db *Database) GetPost(id string) (*Post, error) {
	var post Post
	err := db.DB.First(&post, "id = ?", id).Error
	return &post, err
}

func (db *Database) UpdatePost(post *Post) error {
	return db.DB.Save(post).Error
}

func (db *Database) DeletePost(id string) error {
	return db.DB.Delete(&Post{}, "id = ?", id).Error
}

// Методы для работы с пользователями
func (db *Database) CreateUser(user *User) error {
	return db.DB.Create(user).Error
}

func (db *Database) GetUser(email string) (*User, error) {
	var user User
	err := db.DB.First(&user, "email = ?", email).Error
	return &user, err
}

func (db *Database) UpdateUser(user *User) error {
	return db.DB.Save(user).Error
}

// Методы для работы с токенами
func (db *Database) SaveToken(token *Token) error {
	return db.DB.Create(token).Error
}

func (db *Database) GetToken(refreshToken string) (*Token, error) {
	var token Token
	err := db.DB.First(&token, "refresh_token = ?", refreshToken).Error
	return &token, err
}

func (db *Database) DeleteToken(refreshToken string) error {
	return db.DB.Delete(&Token{}, "refresh_token = ?", refreshToken).Error
}

// Методы для проверки существования
func (db *Database) EmailExists(email string) bool {
	var count int64
	db.DB.Model(&User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func (db *Database) PostExists(id string) bool {
	var count int64
	db.DB.Model(&Post{}).Where("id = ?", id).Count(&count)
	return count > 0
}
