package services

import (
	"api/internal/models"
	"errors"
	"time"
  
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)
type AuthService struct {
	users      map[string]models.User
	tokens     map[string]models.Token
	tokenTTL   time.Duration
	refreshTTL time.Duration

}
func CreatePost(title, content string) (string, error) {
	for _, post := range models.Posts {
		if post.Title == title {
			return "", errors.New("уникальный ключ уже используется")
		}
	}

	postID := uuid.NewString()
	models.Posts[postID] = models.Post{
		ID:      postID,
		Title:   title,
		Content: content,
		Status:  "Draft",
	}

	return postID, nil
}

func GetPublishedPosts() map[string]models.Post {
	publishedPosts := make(map[string]models.Post)

	for id, post := range models.Posts {
		if post.Status == "Published" {
			publishedPosts[id] = post
		}
	}

	return publishedPosts
}

func UpdatePost(postID, title, content string) error {
	post, exists := models.Posts[postID]
	if !exists {
		return errors.New("пост не найден")
	}

	post.Title = title
	post.Content = content
	models.Posts[postID] = post

	return nil
}
func PublishPost(postID, status string) error {
	if status != "Published" {
		return errors.New("статус должен быть Published")
	}

	post, exists := models.Posts[postID]
	if !exists {
		return errors.New("пост не найден")
	}

	post.Status = status
	models.Posts[postID] = post

	return nil
}
func AddImageToPost(postID, imageURL string) error {
	post, exists := models.Posts[postID]
	if !exists {
		return errors.New("пост не найден")
	}

	post.Content += "[Image: " + imageURL + "]"
	models.Posts[postID] = post
	return nil
}

func RemoveImageFromPost(postID string) error {
	post, exists := models.Posts[postID]
	if !exists {
		return errors.New("пост не найден")
	}

	post.Content = "[Image removed]"
	models.Posts[postID] = post
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