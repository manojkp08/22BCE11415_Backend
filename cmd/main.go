package main

import (
	"github.com/gin-gonic/gin"
	"github.com/manojkp08/22BCE11415_Backend/internal/auth"
	"github.com/manojkp08/22BCE11415_Backend/internal/config"
	"github.com/manojkp08/22BCE11415_Backend/internal/database"
	"github.com/manojkp08/22BCE11415_Backend/internal/handlers"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	err := database.InitDB(cfg.DBConnection)
	if err != nil {
		panic(err)
	}

	// Initialize Google OAuth
	auth.InitGoogleOAuth(
		cfg.GoogleClientID,
		cfg.GoogleSecret,
		cfg.GoogleRedirectURL,
	)

	// Set up router
	router := gin.Default()

	// Auth routes
	router.GET("/auth/google/login", handlers.GoogleLoginHandler)
	router.GET("/auth/google/callback", handlers.GoogleCallbackHandler)

	// Start server
	router.Run(":8080")
}
