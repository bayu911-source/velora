
package agents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"velora/internal/services"
)

func TestDataExtractorAgent_Execute(t *testing.T) {
	llm := &MockLLMService{
		GenerateFunc: func(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error) {
			return `{"name": "John Doe", "email": "john.doe@example.com"}`,
			nil
		},
	}

	agent := NewDataExtractorAgent(llm)
	output, err := agent.Execute(context.Background(), "The user's name is John Doe and his email is john.doe@example.com")

	assert.NoError(t, err)
	assert.Equal(t, `{"name": "John Doe", "email": "john.doe@example.com"}`, output)
}
