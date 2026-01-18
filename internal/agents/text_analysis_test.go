
package agents

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextAnalysisAgent_Run(t *testing.T) {
	llm := &MockLLMService{
		GenerateFunc: func(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error) {
			return "test analysis", nil
		},
	}

	agent := NewTextAnalysisAgent(llm)
	output, err := agent.Run("analyze this text")

	assert.NoError(t, err)
	assert.Equal(t, "test analysis", output)
}
