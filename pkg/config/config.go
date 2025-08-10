// backend/pkg/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
}

// Load reads configuration from environment variables.
// It first tries to load a .env file for local development.
func Load() (*Config, error) {
	// Attempt to load .env file. In production, this will likely not exist, which is fine.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/amin_n_co?sslmode=disable"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
	}

	return cfg, nil
}

// getEnv is a helper to read an environment variable or return a default value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
