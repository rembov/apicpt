package main

import (
	"apicpt/internal/models"
	"apicpt/internal/router"
	"apicpt/internal/services"
	"time"
)

func main() {
	db := models.InitDB()

	tokenTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour
	authService := services.NewAuthService(tokenTTL, refreshTTL)

	r := router.SetupRouter(db, authService)
	r.Run(":1488")
}
