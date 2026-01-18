
package agents

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebResearchAgent_Run(t *testing.T) {
	// Test case 1: Successful web research
	t.Run("successful web research", func(t *testing.T) {
		llm := &MockLLMService{
			GenerateFunc: func(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error) {
				return "The capital of France is Paris.", nil
			},
		}

		agent := NewWebResearchAgent(llm)
		input := "what is the capital of france"

		output, err := agent.Run(input)

		assert.NoError(t, err)
		assert.Equal(t, "The capital of France is Paris.", output)
	})

	// Test case 2: LLM service returns an error
	t.Run("llm service error", func(t *testing.T) {
		llm := &MockLLMService{
			GenerateFunc: func(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error) {
				return "", assert.AnError
			},
		}

		agent := NewWebResearchAgent(llm)
		input := "what is the capital of france"

		output, err := agent.Run(input)

		assert.Error(t, err)
		assert.Empty(t, output)
	})
}
