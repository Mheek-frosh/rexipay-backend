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
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect Database
	config.ConnectDB()

	// Set Gin mode (important for production)
	gin.SetMode(gin.ReleaseMode)

	// Initialize router
	r := gin.Default()

	// Health check (ROOT ROUTE)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "RexiPay API is running 🚀",
		})
	})

	// Handle unknown routes (DEBUGGING 🔥)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error": "Route not found",
			"path":  c.Request.URL.Path,
		})
	})

	// Setup all routes
	routes.SetupRoutes(r)

	// Get port from Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port:", port)

	// Run server
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
