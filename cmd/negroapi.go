package main

import (
	"log"

	"api/internal/router"
)

func main() {
	r := router.SetupRouter()
	log.Fatal(r.Run(":1488"))
}
