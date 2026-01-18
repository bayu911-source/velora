
package workflows

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"velora/internal/agents"
)

// MockAgent is a mock implementation of the Agent interface for testing.
type MockAgent struct {
	RunFunc  func(ctx context.Context, input string) (string, error)
	NameFunc func() string
}

// Run calls the mock RunFunc.
func (m *MockAgent) Run(ctx context.Context, input string) (string, error) {
	return m.RunFunc(ctx, input)
}

// Name calls the mock NameFunc.
func (m *MockAgent) Name() string {
	return m.NameFunc()
}

// Description is not used in this test.
func (m *MockAgent) Description() string {
	return ""
}

func TestWorkflow_Run(t *testing.T) {
	// Test case 1: Successful workflow execution
	t.Run("successful workflow execution", func(t *testing.T) {
		agent1 := &MockAgent{
			RunFunc: func(ctx context.Context, input string) (string, error) {
				return "step 1 output", nil
			},
			NameFunc: func() string {
				return "agent1"
			},
		}

		agent2 := &MockAgent{
			RunFunc: func(ctx context.Context, input string) (string, error) {
				assert.Equal(t, "step 1 output", input)
				return "final output", nil
			},
			NameFunc: func() string {
				return "agent2"
			},
		}

		registry := agents.NewRegistry(nil)
		registry.Register(agent1)
		registry.Register(agent2)

		workflow, err := New("test-workflow", []string{"agent1", "agent2"})
		assert.NoError(t, err)

		output, err := workflow.Run(context.Background(), registry, "initial input")

		assert.NoError(t, err)
		assert.Equal(t, "final output", output)
		assert.Equal(t, WorkflowStateCompleted, workflow.State)

		// Verify steps were saved
		steps, err := workflow.GetSteps()
		assert.NoError(t, err)
		assert.Len(t, steps, 2)
		assert.Equal(t, "agent1", steps[0].AgentName)
		assert.Equal(t, "initial input", steps[0].Input)
		assert.Equal(t, "step 1 output", steps[0].Output)
		assert.Equal(t, "agent2", steps[1].AgentName)
		assert.Equal(t, "step 1 output", steps[1].Input)
		assert.Equal(t, "final output", steps[1].Output)
	})

	// Test case 2: Workflow execution fails at a step
	t.Run("workflow execution fails", func(t *testing.T) {
		agent1 := &MockAgent{
			RunFunc: func(ctx context.Context, input string) (string, error) {
				return "step 1 output", nil
			},
			NameFunc: func() string {
				return "agent1"
			},
		}

		expectedErr := errors.New("agent 2 failed")
		agent2 := &MockAgent{
			RunFunc: func(ctx context.Context, input string) (string, error) {
				return "", expectedErr
			},
			NameFunc: func() string {
				return "agent2"
			},
		}

		registry := agents.NewRegistry(nil)
		registry.Register(agent1)
		registry.Register(agent2)

		workflow, err := New("failing-workflow", []string{"agent1", "agent2"})
		assert.NoError(t, err)

		output, err := workflow.Run(context.Background(), registry, "initial input")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedErr.Error())
		assert.Equal(t, "", output)
		assert.Equal(t, WorkflowStateFailed, workflow.State)
	})
}
