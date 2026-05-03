package agents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebResearchAgent_Execute(t *testing.T) {
	// Test case 1: Successful web research
	t.Run("successful web research", func(t *testing.T) {
		mockLLM := &MockLLM{
			GenerateFunc: func(ctx context.Context, prompt string) (string, error) {
				return "The capital of France is Paris.", nil
			},
		}

		agent := NewWebResearchAgent(mockLLM)
		input := "what is the capital of france"

		output, err := agent.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, "The capital of France is Paris.", output)
	})
}
