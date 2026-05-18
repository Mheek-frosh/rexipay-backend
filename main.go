package main

import (
	"log"
	"os"

	"github.com/Mheek-frosh/rexipaybackend/config"
	"github.com/Mheek-frosh/rexipaybackend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Database
	config.ConnectDB()

	// Initialize Gin Router
	r := gin.Default()

	// Setup Routes
	routes.SetupRoutes(r)

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start Server
	r.Run(":" + port)
}
