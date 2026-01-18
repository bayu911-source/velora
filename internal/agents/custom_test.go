
package agents

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomAgent_Run(t *testing.T) {
	llm := &MockLLMService{
		GenerateFunc: func(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error) {
			return "test output", nil
		},
	}

	agent := NewCustomAgent("test-agent", "test prompt", llm)
	output, err := agent.Run("test input")

	assert.NoError(t, err)
	assert.Equal(t, "test output", output)
}
