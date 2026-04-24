
package config

import (
	"os"
)

// Config stores the application configuration.
type Config struct {
	GeminiAPIKey   string `mapstructure:"GEMINI_API_KEY"`
	GeminiAPIURL   string `mapstructure:"GEMINI_API_URL"`
	OpenAIAPIKey   string `mapstructure:"OPENAI_API_KEY"`
}

// LoadConfig loads the configuration from environment variables.
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		GeminiAPIURL: os.Getenv("GEMINI_API_URL"),
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
	}

	// Set default Gemini URL if not set
	if cfg.GeminiAPIURL == "" {
		cfg.GeminiAPIURL = "https://generativelanguage.googleapis.com"
	}

	return cfg, nil
}
