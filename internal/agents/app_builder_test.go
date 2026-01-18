
package agents

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"velora/internal/services"
)

// MockLLMService is a mock implementation of the LLMService interface for testing.
type MockLLMService struct {
	GenerateFunc func(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error)
    GenerateStreamFunc func(prompt, modelName string, temperature float32, maxOutputTokens int) (<-chan string, error)
}

// Generate calls the mock GenerateFunc.
func (m *MockLLMService) Generate(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error) {
	return m.GenerateFunc(prompt, modelName, temperature, maxOutputTokens)
}

// GenerateStream calls the mock GenerateStreamFunc.
func (m *MockLLMService) GenerateStream(prompt, modelName string, temperature float32, maxOutputTokens int) (<-chan string, error) {
	return m.GenerateStreamFunc(prompt, modelName, temperature, maxOutputTokens)
}

func TestAppBuilderAgent_Run(t *testing.T) {
	// Test case 1: Successful code generation
	t.Run("successful code generation", func(t *testing.T) {
		llm := &MockLLMService{
			GenerateFunc: func(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error) {
				return "```go\nfmt.Println(\"Hello, World!\")\n```", nil
			},
		}

		agent := NewAppBuilderAgent(llm)
		input := "write a hello world program in go"

		output, err := agent.Run(input)

		assert.NoError(t, err)
		assert.Equal(t, "```go\nfmt.Println(\"Hello, World!\")\n```", output)
	})

	// Test case 2: LLM service returns an error
	t.Run("llm service error", func(t *testing.T) {
		llm := &MockLLMService{
			GenerateFunc: func(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error) {
				return "", assert.AnError
			},
		}

		agent := NewAppBuilderAgent(llm)
		input := "write a hello world program in go"

		output, err := agent.Run(input)

		assert.Error(t, err)
		assert.Empty(t, output)
	})
}
