package main

import (
	"log"
	"ws/internal/handlers"
)

func main() {

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	log.Println("Listening on :8080")
	err := routes().Run(":8080")
	if err != nil {
		log.Fatal(err)
	}

}
