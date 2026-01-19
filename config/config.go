
package config

// Config stores the application configuration.
type Config struct {
	GeminiAPIKey   string `mapstructure:"GEMINI_API_KEY"`
	GeminiAPIURL   string `mapstructure:"GEMINI_API_URL"`
	OpenAIAPIKey string `mapstructure:"OPENAI_API_KEY"`
}

// LoadConfig loads the configuration from a file.
func LoadConfig(path string) (*Config, error) {
	// In a real application, you would load the configuration from a file.
	// For now, we'll just return a default config.
	return &Config{}, nil
}
