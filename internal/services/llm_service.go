
package services

import (
	"fmt"

	"velora/config"
)

// LLMService is a generic interface for a language model service.
type LLMService interface {
	Generate(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error)
	GenerateStream(prompt, modelName string, temperature float32, maxOutputTokens int) (<-chan string, <-chan error)
	Close() error
}

// NewLLMService creates a new LLMService based on the provided configuration.
func NewLLMService(cfg config.Config) (LLMService, error) {
	if cfg.GeminiAPIKey != "" {
		return NewGeminiService(cfg)
	}

	if cfg.OpenAIAPIKey != "" {
		return NewOpenAIService(cfg)
	}

	return nil, fmt.Errorf("no LLM service configured")
}
