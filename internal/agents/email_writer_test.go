
package agents

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"velora/internal/services"
)

// MockLLM is a mock implementation of the LLM interface for testing.
type MockLLM struct {
	GenerateFunc func(ctx context.Context, prompt string) (string, error)
}

// Generate calls the mock GenerateFunc.
func (m *MockLLM) Generate(ctx context.Context, prompt string) (string, error) {
	return m.GenerateFunc(ctx, prompt)
}

// Close is a no-op for the mock.
func (m *MockLLM) Close() {}

func TestEmailWriterAgent_Execute(t *testing.T) {
	llm := &MockLLM{
		GenerateFunc: func(ctx context.Context, prompt string) (string, error) {
			return "test email", nil
		},
	}

	agent := NewEmailWriterAgent(llm)
	output, err := agent.Execute(context.Background(), "write a test email")

	assert.NoError(t, err)
	assert.Equal(t, "test email", output)
}
