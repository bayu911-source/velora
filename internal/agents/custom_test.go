
package agents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomAgent_Execute(t *testing.T) {
	mockLLM := &MockLLM{
		GenerateFunc: func(ctx context.Context, prompt string) (string, error) {
			return "test output", nil
		},
	}

	agent := NewCustomAgent("test-agent", "test description", "test prompt", mockLLM)
	output, err := agent.Execute(context.Background(), "test input")

	assert.NoError(t, err)
	assert.Equal(t, "test output", output)
}
