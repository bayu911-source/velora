package agents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextAnalysisAgent_Execute(t *testing.T) {
	mockLLM := &MockLLM{
		GenerateFunc: func(ctx context.Context, prompt string) (string, error) {
			return "test analysis", nil
		},
	}

	agent := NewTextAnalysisAgent(mockLLM)
	output, err := agent.Execute(context.Background(), "analyze this text")

	assert.NoError(t, err)
	assert.Equal(t, "test analysis", output)
}
