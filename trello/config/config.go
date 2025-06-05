package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config struct holds all configuration for the application
type Config struct {
	DBHost       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBPort       string
	DBSSLMode    string
	DBTimeZone   string
	JWTSecretKey string
	ServerPort   string
}

// LoadConfig loads configuration from .env file or environment variables
func LoadConfig() (*Config, error) {
	// Attempt to load .env file, but don't fail if it's not present
	// This allows for environments where vars are set directly (e.g., Docker, K8s)
	_ = godotenv.Load(".env.example")

	return &Config{
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "password"),
		DBName:       getEnv("DB_NAME", "trello_clone_db"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBSSLMode:    getEnv("DB_SSLMODE", "disable"),
		DBTimeZone:   getEnv("DB_TIMEZONE", "UTC"),
		JWTSecretKey: getEnv("JWT_SECRET_KEY", "defaultsecret"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
	}, nil
}

// Helper function to get an environment variable or return a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
