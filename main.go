package main

import (
	"log"
	"time"

	"contrihub/database"
	"contrihub/models"
	"contrihub/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load. Using system environment variables.")
	}

	// Initialize database connection
	database.Connect()

	// Initialize Redis connection (for rate limiting)
	database.ConnectRedis()

	// Run auto-migrations
	if err := database.DB.AutoMigrate(&models.AICache{}, &models.AIChatMessage{}); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Println("Database migrations completed.")

	// Start background cleanup for chats older than 30 days
	go func() {
		for {
			database.DB.Where("created_at < now() - interval '30 days'").Delete(&models.AIChatMessage{})
			time.Sleep(24 * time.Hour)
		}
	}()

	r := router.SetupRouter()

	log.Println("Server starting on port 5050...")
	if err := r.Run(":5050"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
