package handlers

import (
	"api/internal/middleware"
	"api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService *services.AuthService
}

func NewHandler(authService *services.AuthService) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) RegisterHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.RegisterUser(input.Email, input.Password, input.Role)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Регистрация прошла успешно"})
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.authService.AuthenticateUser(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := middleware.GenerateToken(input.Email, role, h.authService.GetTokenTTL())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания токена"})
		return
	}

	refreshToken := h.authService.GenerateRefreshToken(input.Email)

	c.SetCookie("accessToken", accessToken, int(h.authService.GetTokenTTL().Seconds()), "/", "", false, true)
	c.SetCookie("refreshToken", refreshToken, int(h.authService.GetRefreshTTL().Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Вход выполнен успешно"})
}

func (h *Handler) RefreshTokenHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh токен не найден"})
		return
	}

	email, err := h.authService.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, _ := h.authService.GetUserRole(email)
	newAccessToken, err := middleware.GenerateToken(email, role, h.authService.GetTokenTTL())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	c.SetCookie("accessToken", newAccessToken, int(h.authService.GetTokenTTL().Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Токен успешно обновлён"})
}
