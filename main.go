package main

import (
	"api/internal/models"
	"api/internal/router"
	"api/internal/services"
	"time"
)

func main() {
	// Инициализация базы данных
	db := models.InitDB()

	// Установите время жизни токенов
	tokenTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour

	// Создание AuthService
	authService := services.NewAuthService(tokenTTL, refreshTTL)

	// Создание роутера
	r := router.SetupRouter(db, authService)
	r.Run(":1488")
}

