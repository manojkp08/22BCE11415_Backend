package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GoogleClientID    string
	GoogleSecret      string
	GoogleRedirectURL string
	DBConnection      string
	JWTSecret         string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	return &Config{
		GoogleClientID:    getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleSecret:      getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL: getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		DBConnection:      getEnv("DB_CONNECTION", "postgres://user:password@localhost:5432/file_sharing?sslmode=disable"),
		JWTSecret:         getEnv("JWT_SECRET", "your_secret_key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
