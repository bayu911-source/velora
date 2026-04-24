
package agents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"velora/internal/services"
)

// MockLLM is a mock implementation of the LLM interface for testing.
type MockLLM struct {
	GenerateFunc func(ctx context.Context, prompt string) (string, error)
}

// Generate calls the mock GenerateFunc.
func (m *MockLLM) Generate(ctx context.Context, prompt string) (string, error) {
	return m.GenerateFunc(ctx, prompt)
}

// Close is a no-op for the mock.
func (m *MockLLM) Close() {}

func TestWebResearchAgent_Execute(t *testing.T) {
	// Test case 1: Successful web research
	t.Run("successful web research", func(t *testing.T) {
		llm := &MockLLM{
			GenerateFunc: func(ctx context.Context, prompt string) (string, error) {
				return "The capital of France is Paris.", nil
			},
		}

		agent := NewWebResearchAgent(llm)
		input := "what is the capital of france"

		output, err := agent.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, "The capital of France is Paris.", output)
	})

	// Test case 2: LLM service returns an error
	t.Run("llm service error", func(t *testing.T) {
		llm := &MockLLM{
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
