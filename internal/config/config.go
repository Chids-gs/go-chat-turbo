package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type Config struct {
	GeminiAPIKey string
	GeminiAPIURL string
}

func LoadConfig() *Config {
	// Get the project root directory
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "..", "..")

	// Load .env from project root
	if err := godotenv.Load(filepath.Join(projectRoot, ".env")); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	config := &Config{
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		GeminiAPIURL: os.Getenv("GEMINI_API_URL"),
	}

	// Validate required configuration
	if config.GeminiAPIKey == "" {
		log.Fatal("GEMINI_API_KEY is required")
	}

	if config.GeminiAPIURL == "" {
		log.Fatal("GEMINI_API_URL is required")
	}

	return config
}
