package main

import (
	"log"
	"ws/internal/handlers"
)

func main() {
	err := routes().Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on :8080")
	log.Println("Starting channel listener")
	go handlers.ListenForWsChan()
}
