package agents

import (
	"context"
)

// MockLLM is a mock implementation of LLMService for testing.
type MockLLM struct {
	GenerateFunc func(ctx context.Context, prompt string) (string, error)
}

// Generate calls the mock GenerateFunc.
func (m *MockLLM) Generate(ctx context.Context, prompt string) (string, error) {
	if m.GenerateFunc != nil {
		return m.GenerateFunc(ctx, prompt)
	}
	return "mock response", nil
}

// Close is a no-op for the mock.
func (m *MockLLM) Close() {}
