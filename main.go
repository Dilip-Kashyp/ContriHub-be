package main

import (
	"log"

	"contrihub/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load. Using system environment variables.")
	}

	r := router.SetupRouter()

	log.Println("Server starting on port 5050...")
	if err := r.Run(":5050"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
