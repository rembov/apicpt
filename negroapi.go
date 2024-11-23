package main

import (
	"log"

	"api/router"
)

func main() {
	r := router.SetupRouter()
	log.Fatal(r.Run(":1488"))
}
