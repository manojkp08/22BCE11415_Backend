package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manojkp08/22BCE11415_Backend/internal/auth"
	"github.com/manojkp08/22BCE11415_Backend/internal/cache"
	"github.com/manojkp08/22BCE11415_Backend/internal/config"
	"github.com/manojkp08/22BCE11415_Backend/internal/database"
	"github.com/manojkp08/22BCE11415_Backend/internal/handlers"
	"github.com/manojkp08/22BCE11415_Backend/internal/worker"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	go worker.StartCleanupWorker(24 * time.Hour) // Runs daily

	// Initialize database
	err := database.InitDB(cfg.DBConnection)
	if err != nil {
		panic(err)
	}

	// Initialize Redis cache
	if err := cache.InitRedis(cfg.RedisAddr, cfg.RedisPassword); err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}

	// Initialize Google OAuth
	auth.InitGoogleOAuth(
		cfg.GoogleClientID,
		cfg.GoogleSecret,
		cfg.GoogleRedirectURL,
	)

	// Set up router
	router := gin.Default()
	handlers.SetupRoutes(router)

	router.Run(":8080")
}
