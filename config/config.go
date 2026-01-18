
package config

import (
	"github.com/spf13/viper"
)

// Config stores all configuration for the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	GeminiAPIKey    string            `mapstructure:"GEMINI_API_KEY"`
	OpenAIAPIKey    string            `mapstructure:"OPENAI_API_KEY"`
	LLMProvider     string            `mapstructure:"LLM_PROVIDER"`
	GeminiAPIURL    string            `mapstructure:"GEMINI_API_URL"`
	CustomModels    map[string]string `mapstructure:"CUSTOM_MODELS"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
