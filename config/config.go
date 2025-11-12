package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	// Middleware API Configuration
	MiddlewareAPIKey  string
	MiddlewareBaseURL string

	// Application Mode: stdio, http, sse
	AppMode string

	// Server Configuration (for http/sse modes)
	AppHost string
	AppPort string

	// Tool Exclusion
	ExcludedTools map[string]bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	cfg := &Config{
		MiddlewareAPIKey:  os.Getenv("MIDDLEWARE_API_KEY"),
		MiddlewareBaseURL: os.Getenv("MIDDLEWARE_BASE_URL"),
		AppMode:           getEnvOrDefault("APP_MODE", "stdio"),
		AppHost:           getEnvOrDefault("APP_HOST", "localhost"),
		AppPort:           getEnvOrDefault("APP_PORT", "8080"),
		ExcludedTools:     make(map[string]bool),
	}

	if cfg.MiddlewareAPIKey == "" {
		return nil, fmt.Errorf("MIDDLEWARE_API_KEY is required")
	}
	if cfg.MiddlewareBaseURL == "" {
		return nil, fmt.Errorf("MIDDLEWARE_BASE_URL is required")
	}

	if excludedStr := os.Getenv("EXCLUDED_TOOLS"); excludedStr != "" {
		for _, tool := range strings.Split(excludedStr, ",") {
			tool = strings.TrimSpace(tool)
			if tool != "" {
				cfg.ExcludedTools[tool] = true
			}
		}
	}

	validModes := map[string]bool{"stdio": true, "http": true, "sse": true}
	if !validModes[cfg.AppMode] {
		return nil, fmt.Errorf("invalid APP_MODE: %s (must be stdio, http, or sse)", cfg.AppMode)
	}

	return cfg, nil
}

// IsToolExcluded checks if a tool is excluded
func (c *Config) IsToolExcluded(toolName string) bool {
	return c.ExcludedTools[toolName]
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

