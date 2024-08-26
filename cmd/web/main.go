package main

import "log"

func main() {
	err := routes().Run("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
}
