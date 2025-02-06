package main

import (
	"apicpt/internal/api"
	"apicpt/internal/entites"
	"apicpt/internal/services"
	"time"
)

func main() {
	db := models.InitDB()

	tokenTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour
	authService := services.NewAuthService(tokenTTL, refreshTTL)

	r := api.SetupRouter(db, authService)
	r.Run(":1488") // можно сделать файл с конфигом
}
