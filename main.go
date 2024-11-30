package main

import (
	"log"

	"api/internal/router"
	"api/internal/services"
)

func main() {
	authService := &services.AuthService{}
	r := router.SetupRouter(authService)
	log.Fatal(r.Run(":1488"))
}
