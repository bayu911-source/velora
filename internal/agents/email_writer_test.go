
package agents

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailWriterAgent_Run(t *testing.T) {
	llm := &MockLLMService{
		GenerateFunc: func(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error) {
			return "test email", nil
		},
	}

	agent := NewEmailWriterAgent(llm)
	output, err := agent.Run("write a test email")

	assert.NoError(t, err)
	assert.Equal(t, "test email", output)
}
