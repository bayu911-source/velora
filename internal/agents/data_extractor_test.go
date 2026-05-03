
package agents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataExtractorAgent_Execute(t *testing.T) {
	mockLLM := &MockLLM{
		GenerateFunc: func(ctx context.Context, prompt string) (string, error) {
			return `{"name": "John Doe", "email": "john.doe@example.com"}`, nil
		},
	}

	agent := NewDataExtractorAgent(mockLLM)
	output, err := agent.Execute(context.Background(), "The user's name is John Doe and his email is john.doe@example.com")

	assert.NoError(t, err)
	assert.Equal(t, `{"name": "John Doe", "email": "john.doe@example.com"}`, output)
}
