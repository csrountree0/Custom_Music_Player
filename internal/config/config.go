package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database connection fields
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	// HTTP server
	ServerPort string

	// JWT authentication
	JWTSecret          string
	JWTExpiry          time.Duration
	RefreshTokenExpiry time.Duration

	// File storage
	AudioDir string

	// Runtime environment
	Environment string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBName:      getEnv("DB_NAME", "musicapp"),
		DBUser:      getEnv("DB_USER", "musicapp"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		AudioDir:    getEnv("AUDIO_DIR", "./audio"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	var err error
	cfg.JWTExpiry, err = time.ParseDuration(getEnv("JWT_EXPIRY", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY: %w", err)
	}

	cfg.RefreshTokenExpiry, err = time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_EXPIRY: %w", err)
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
