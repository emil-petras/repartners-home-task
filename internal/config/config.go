package config

import (
	"fmt"
	"os"
)

// Config holds application configuration settings
type Config struct {
	ServerPort  string // port for the http server
	DatabaseURL string // path to the sqlite database file
}

// Load creates a new config with environment variables or defaults
func Load() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "./data/pack_sizes.db"),
	}
}

// GetServerAddr returns the full server address string
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf(":%s", c.ServerPort)
}

// getEnv retrieves environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
