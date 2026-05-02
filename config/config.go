
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config stores the application configuration.
type Config struct {
	AppName         string        `mapstructure:"APP_NAME"`
	Environment     string        `mapstructure:"APP_ENV"`
	Port            string        `mapstructure:"APP_PORT"`
	DatabaseURL     string        `mapstructure:"DATABASE_URL"`
	RedisURL        string        `mapstructure:"REDIS_URL"`
	JWTSecret       string        `mapstructure:"JWT_SECRET"`
	RefreshSecret   string        `mapstructure:"JWT_REFRESH_SECRET"`
	GeminiAPIKey    string        `mapstructure:"GEMINI_API_KEY"`
	GeminiAPIURL    string        `mapstructure:"GEMINI_API_URL"`
	OpenAIAPIKey    string        `mapstructure:"OPENAI_API_KEY"`
	RateLimitMax    int           `mapstructure:"RATE_LIMIT_MAX"`
	RateLimitWindow time.Duration `mapstructure:"RATE_LIMIT_WINDOW"`
	AccessTokenTTL  time.Duration `mapstructure:"ACCESS_TOKEN_TTL"`
	RefreshTokenTTL time.Duration `mapstructure:"REFRESH_TOKEN_TTL"`
}

// LoadConfig loads the configuration from environment variables and .env file.
func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("env")
	v.SetConfigName(".env")
	v.AddConfigPath(path)
	v.AutomaticEnv()

	v.SetDefault("APP_NAME", "Velora AI Automation Framework")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("REDIS_URL", "redis:6379")
	v.SetDefault("RATE_LIMIT_MAX", 500)
	v.SetDefault("RATE_LIMIT_WINDOW", "1m")
	v.SetDefault("ACCESS_TOKEN_TTL", "15m")
	v.SetDefault("REFRESH_TOKEN_TTL", "168h")
	v.SetDefault("GEMINI_API_URL", "https://generativelanguage.googleapis.com")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	cfg := &Config{
		AppName:         v.GetString("APP_NAME"),
		Environment:     v.GetString("APP_ENV"),
		Port:            v.GetString("APP_PORT"),
		DatabaseURL:     v.GetString("DATABASE_URL"),
		RedisURL:        v.GetString("REDIS_URL"),
		JWTSecret:       v.GetString("JWT_SECRET"),
		RefreshSecret:   v.GetString("JWT_REFRESH_SECRET"),
		GeminiAPIKey:    v.GetString("GEMINI_API_KEY"),
		GeminiAPIURL:    v.GetString("GEMINI_API_URL"),
		OpenAIAPIKey:    v.GetString("OPENAI_API_KEY"),
		RateLimitMax:    v.GetInt("RATE_LIMIT_MAX"),
		RateLimitWindow: v.GetDuration("RATE_LIMIT_WINDOW"),
		AccessTokenTTL:  v.GetDuration("ACCESS_TOKEN_TTL"),
		RefreshTokenTTL: v.GetDuration("REFRESH_TOKEN_TTL"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.RedisURL == "" {
		return nil, fmt.Errorf("REDIS_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.RefreshSecret == "" {
		return nil, fmt.Errorf("JWT_REFRESH_SECRET is required")
	}

	if cfg.RateLimitWindow == 0 {
		cfg.RateLimitWindow = time.Minute
	}
	if cfg.AccessTokenTTL == 0 {
		cfg.AccessTokenTTL = 15 * time.Minute
	}
	if cfg.RefreshTokenTTL == 0 {
		cfg.RefreshTokenTTL = 168 * time.Hour
	}

	return cfg, nil
}
