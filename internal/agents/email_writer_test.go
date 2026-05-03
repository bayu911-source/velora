
package agents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailWriterAgent_Execute(t *testing.T) {
	mockLLM := &MockLLM{
		GenerateFunc: func(ctx context.Context, prompt string) (string, error) {
			return "test email", nil
		},
	}

	agent := NewEmailWriterAgent(mockLLM)
	output, err := agent.Execute(context.Background(), "write a test email")

	assert.NoError(t, err)
	assert.Equal(t, "test email", output)
}
