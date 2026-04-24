
package agents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"velora/internal/services"
)

func TestCustomAgent_Execute(t *testing.T) {
	llm := &MockLLMService{
		GenerateFunc: func(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error) {
			return "test output", nil
		},
	}

	agent := NewCustomAgent("test-agent", "test description", "test prompt", llm)
	output, err := agent.Execute(context.Background(), "test input")

	assert.NoError(t, err)
	assert.Equal(t, "test output", output)
}
