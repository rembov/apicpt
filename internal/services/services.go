package services

import (
	"apicpt/internal/entites"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	users      map[string]models.User
	tokens     map[string]models.Token
	tokenTTL   time.Duration
	refreshTTL time.Duration
}

func CreatePost(db *gorm.DB, title, content string, authorID uint) (string, error) {
	// Проверка уникальности заголовка
	var count int64
	db.Model(&models.Post{}).Where("title = ?", title).Count(&count)
	if count > 0 {
		return "", errors.New("уникальный ключ уже используется")
	}

	// Создание нового поста
	post := models.Post{
		ID:       uuid.NewString(),
		Title:    title,
		Content:  content,
		Status:   "Draft",
		AuthorID: authorID,
	}
	if err := db.Create(&post).Error; err != nil {
		return "", err
	}

	return post.ID, nil
}

func GetPublishedPosts(db *gorm.DB) ([]models.Post, error) {
	var posts []models.Post
	if err := db.Where("status = ?", "Published").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func UpdatePost(db *gorm.DB, postID, title, content string) error {
	var post models.Post
	if err := db.First(&post, "id = ?", postID).Error; err != nil {
		return errors.New("пост не найден")
	}

	post.Title = title
	post.Content = content
	if err := db.Save(&post).Error; err != nil {
		return err
	}

	return nil
}
func PublishPost(db *gorm.DB, postID, status string) error {
	if status != "Published" {
		return errors.New("статус должен быть Published")
	}

	var post models.Post
	if err := db.First(&post, "id = ?", postID).Error; err != nil {
		return errors.New("пост не найден")
	}

	post.Status = status
	if err := db.Save(&post).Error; err != nil {
		return err
	}

	return nil
}

func AddImageToPost(db *gorm.DB, postID, imageURL string) error {
	var post models.Post
	if err := db.First(&post, "id = ?", postID).Error; err != nil {
		return errors.New("пост не найден")
	}

	// Добавляем информацию об изображении
	post.Content += "\n[Image: " + imageURL + "]"
	if err := db.Save(&post).Error; err != nil {
		return err
	}

	return nil
}
func RemoveImageFromPost(db *gorm.DB, postID, imageID string) error {
	var post models.Post
	if err := db.First(&post, "id = ?", postID).Error; err != nil {
		return errors.New("пост не найден")
	}

	// Удаляем изображение
	post.Content = "[Image removed]"
	if err := db.Save(&post).Error; err != nil {
		return err
	}

	return nil
}
func NewAuthService(tokenTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		users:      make(map[string]models.User),
		tokens:     make(map[string]models.Token),
		tokenTTL:   tokenTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *AuthService) RegisterUser(email, password, role string) error {
	if _, exists := s.users[email]; exists {
		return errors.New("email уже зарегистрирован")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	s.users[email] = models.User{
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}
	return nil
}

func (s *AuthService) AuthenticateUser(email, password string) (string, error) {
	user, exists := s.users[email]
	if !exists || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("неверный логин или пароль")
	}
	return user.Role, nil
}

func (s *AuthService) GenerateRefreshToken(email string) string {
	refreshToken := uuid.NewString()
	s.tokens[email] = models.Token{
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTTL),
	}
	return refreshToken
}

func (s *AuthService) ValidateRefreshToken(refreshToken string) (string, error) {
	for email, token := range s.tokens {
		if token.RefreshToken == refreshToken && token.ExpiresAt.After(time.Now()) {
			return email, nil
		}
	}
	return "", errors.New("refresh токен недействителен")
}

func (s *AuthService) GetUserRole(email string) (string, error) {
	user, exists := s.users[email]
	if !exists {
		return "", errors.New("пользователь не найден")
	}
	return user.Role, nil
}

func (s *AuthService) GetTokenTTL() time.Duration {
	return s.tokenTTL
}

func (s *AuthService) GetRefreshTTL() time.Duration {
	return s.refreshTTL
}
