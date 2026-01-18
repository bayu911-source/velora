
package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"velora/config"
)

func TestNewGeminiService(t *testing.T) {
	// Test case 1: Valid config
	t.Run("valid config", func(t *testing.T) {
		cfg := config.Config{
			GeminiAPIKey: "test-api-key",
			GeminiAPIURL: "test-url",
		}

		gs, err := NewGeminiService(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, gs)
	})

	// Test case 2: Empty API key
	t.Run("empty api key", func(t *testing.T) {
		cfg := config.Config{
			GeminiAPIKey: "",
			GeminiAPIURL: "test-url",
		}

		gs, err := NewGeminiService(cfg)

		assert.Error(t, err)
		assert.Nil(t, gs)
		assert.Equal(t, "GEMINI_API_KEY environment variable not set", err.Error())
	})
}
