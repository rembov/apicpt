package main

import (
	"api/internal/router"
	"api/internal/models"
)

func main() {
	// Инициализация базы данных
	db := models.InitDB()

	// Создание роутера и передача базы данных
	r := router.SetupRouter(db)
	r.Run(":1488")
}
