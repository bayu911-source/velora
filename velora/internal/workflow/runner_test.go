package workflow

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/velora-id/velora/pkg"
)

// MockAgent is a mock implementation of the Agent interface for testing.
// It allows us to control its output and check what input it received.	ype MockAgent struct {
	// The name of the agent.
	AgentName string
	// The output that the Run method should return.
	RunOutput string
	// Any error that the Run method should return.
	RunError error
	// Records the last input received by the Run method.
	LastInput string
}

func (a *MockAgent) Name() string {
	return a.AgentName
}

// Run records the input and returns the pre-configured output and error.
// It also demonstrates writing to memory.
func (a *MockAgent) Run(memory *pkg.MemoryManager, input string) (string, error) {
	a.LastInput = input
	// Mock writing to memory
	memory.Set(fmt.Sprintf("%s_output", a.AgentName), a.RunOutput)
	return a.RunOutput, a.RunError
}

// TestRunner_Run executes a test on the pipeline runner.
func TestRunner_Run(t *testing.T) {
	// 1. Setup
	runner := NewRunner()

	agent1 := &MockAgent{
		AgentName: "agent1",
		RunOutput: "output from agent1",
	}
	agent2 := &MockAgent{
		AgentName: "agent2",
		RunOutput: "output from agent2",
	}

	runner.RegisterAgent(agent1)
	runner.RegisterAgent(agent2)

	pipeline := &Pipeline{
		Name: "Test Pipeline",
		Steps: []Step{
			{
				Agent: "agent1",
				Input: "initial input",
			},
			{
				Agent: "agent2",
				// Input is implicitly the output of agent1
			},
		},
	}

	// 2. Execute
	finalResult, err := runner.Run(pipeline)

	// 3. Assert
	assert.NoError(t, err) // We expect no error
	assert.Equal(t, "output from agent2", finalResult) // The final result should be from the last agent
	assert.Equal(t, "initial input", agent1.LastInput) // Check input for the first agent
	assert.Equal(t, "output from agent1", agent2.LastInput) // Check that the second agent received the first agent's output
}
