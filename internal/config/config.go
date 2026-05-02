package config

import (
    "fmt"
    "strings"
    "time"

    "github.com/spf13/viper"
)

// Config stores application configuration loaded from environment variables.
type Config struct {
    AppName          string        `mapstructure:"APP_NAME"`
    Env              string        `mapstructure:"APP_ENV"`
    Port             string        `mapstructure:"APP_PORT"`
    DatabaseURL      string        `mapstructure:"DATABASE_URL"`
    RedisURL         string        `mapstructure:"REDIS_URL"`
    JWTSecret        string        `mapstructure:"JWT_SECRET"`
    RefreshSecret    string        `mapstructure:"JWT_REFRESH_SECRET"`
    GeminiAPIKey     string        `mapstructure:"GEMINI_API_KEY"`
    GeminiAPIURL     string        `mapstructure:"GEMINI_API_URL"`
    RateLimitMax     int           `mapstructure:"RATE_LIMIT_MAX"`
    RateLimitWindow  time.Duration `mapstructure:"RATE_LIMIT_WINDOW"`
    AccessTokenTTL   time.Duration
    RefreshTokenTTL  time.Duration
}

// LoadConfig reads environment variables and returns a parsed Config.
func LoadConfig(path string) (*Config, error) {
    v := viper.New()
    v.SetConfigFile(".env")
    v.SetConfigType("env")
    v.AutomaticEnv()
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

    if err := v.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }

    v.SetDefault("APP_NAME", "Velora AI Automation Framework")
    v.SetDefault("APP_ENV", "development")
    v.SetDefault("APP_PORT", "8080")
    v.SetDefault("GEMINI_API_URL", "https://generativelanguage.googleapis.com")
    v.SetDefault("RATE_LIMIT_MAX", 500)
    v.SetDefault("RATE_LIMIT_WINDOW", 1*time.Minute)

    cfg := &Config{
        AppName:         v.GetString("APP_NAME"),
        Env:             v.GetString("APP_ENV"),
        Port:            v.GetString("APP_PORT"),
        DatabaseURL:     v.GetString("DATABASE_URL"),
        RedisURL:        v.GetString("REDIS_URL"),
        JWTSecret:       v.GetString("JWT_SECRET"),
        RefreshSecret:   v.GetString("JWT_REFRESH_SECRET"),
        GeminiAPIKey:    v.GetString("GEMINI_API_KEY"),
        GeminiAPIURL:    v.GetString("GEMINI_API_URL"),
        RateLimitMax:    v.GetInt("RATE_LIMIT_MAX"),
        RateLimitWindow: v.GetDuration("RATE_LIMIT_WINDOW"),
        AccessTokenTTL:  15 * time.Minute,
        RefreshTokenTTL: 7 * 24 * time.Hour,
    }

    missing := []string{}
    if cfg.DatabaseURL == "" {
        missing = append(missing, "DATABASE_URL")
    }
    if cfg.RedisURL == "" {
        missing = append(missing, "REDIS_URL")
    }
    if cfg.JWTSecret == "" {
        missing = append(missing, "JWT_SECRET")
    }
    if cfg.RefreshSecret == "" {
        missing = append(missing, "JWT_REFRESH_SECRET")
    }
    if cfg.GeminiAPIKey == "" {
        missing = append(missing, "GEMINI_API_KEY")
    }
    if len(missing) > 0 {
        return nil, fmt.Errorf("missing required configuration: %s", strings.Join(missing, ", "))
    }

    return cfg, nil
}
