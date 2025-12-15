package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration.
type Config struct {
	Port              string
	OpenWeatherAPIKey string
	OpenAIAPIKey      string
}

// Load reads configuration from environment variables.
// Returns an error if required environment variables are missing.
func Load() (*Config, error) {
	openWeatherKey := os.Getenv("OPENWEATHER_API_KEY")
	if openWeatherKey == "" {
		return nil, fmt.Errorf("OPENWEATHER_API_KEY environment variable is required")
	}

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:              port,
		OpenWeatherAPIKey: openWeatherKey,
		OpenAIAPIKey:      openAIKey,
	}, nil
}
