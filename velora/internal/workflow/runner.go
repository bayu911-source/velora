
package workflow

import (
	"fmt"
	"github.com/velora-id/velora/pkg"
)

type Runner struct {
	agents map[string]pkg.Agent
}

func NewRunner() *Runner {
	return &Runner{
		agents: make(map[string]pkg.Agent),
	}
}

func (r *Runner) RegisterAgent(agent pkg.Agent) {
	r.agents[agent.Name()] = agent
}

func (r *Runner) Run(pipeline *Pipeline) (string, error) {
	var result string
	var err error

	// Create a new MemoryManager for this pipeline run
	memory := pkg.NewMemoryManager()

	for i, step := range pipeline.Steps {
		agent, ok := r.agents[step.Agent]
		if !ok {
			return "", fmt.Errorf("agent '%s' not found", step.Agent)
		}

		input := step.Input
		// On the first step, we can use the pipeline's global input.
		// On subsequent steps, we use the result of the previous step.
		if i > 0 && result != "" {
			input = result
		} else {
			// Also store the initial input in memory so other agents can access it
			memory.Set("initial_input", input)
		}

		// Run the agent with memory and the current input
		result, err = agent.Run(memory, input)
		if err != nil {
			return "", fmt.Errorf("error in agent '%s': %w", step.Agent, err)
		}

		// Store the output of this step in memory for subsequent agents
		memoryKey := fmt.Sprintf("step_%d_%s_output", i, step.Agent)
		memory.Set(memoryKey, result)
	}

	return result, nil
}
