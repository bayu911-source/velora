
package agents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppBuilderAgent_Run(t *testing.T) {
	// Test case 1: Successful code generation
	t.Run("successful code generation", func(t *testing.T) {
		mockLLM := &MockLLM{
			GenerateFunc: func(ctx context.Context, prompt string) (string, error) {
				// Return a properly formatted response with file separator ---
				return "main.go\npackage main\n", nil
			},
		}

		agent := NewAppBuilderAgent(mockLLM)
		input := "write a hello world program in go"

		output, err := agent.Execute(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, "Application built successfully!", output)
	})
}
