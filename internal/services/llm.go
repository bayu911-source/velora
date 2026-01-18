
package services

import (
	"fmt"

	"velora/config"
)

// LLMService defines the interface for a Large Language Model service.
// This allows for mocking in tests and abstracting the underlying implementation.
type LLMService interface {
	Generate(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error)
	GenerateStream(prompt, modelName string, temperature float32, maxOutputTokens int) (<-chan string, <-chan error)
	Close() error
}

// New creates a new LLMService based on the provided configuration.
func New(cfg config.Config) (LLMService, error) {
	switch cfg.LLMProvider {
	case "gemini":
		return NewGeminiService(cfg)
	case "openai":
		return NewOpenAIService(cfg)
	default:
		return nil, fmt.Errorf("unknown LLM provider: %s", cfg.LLMProvider)
	}
}
