
package workflow

import (
	"encoding/json"
	"fmt"
	"strings"
	"velora/internal/agents"
)

// Step defines a single step in a workflow.
type Step struct {
	Agent string `json:"agent"`
	Input string `json:"input"`
}

// Workflow defines a series of steps to be executed.
type Workflow struct {
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
}

// PipelineRunner runs a workflow.
type PipelineRunner struct {
	agents map[string]agents.Agent
}

// NewPipelineRunner creates a new PipelineRunner.
func NewPipelineRunner(agents map[string]agents.Agent) *PipelineRunner {
	return &PipelineRunner{
		agents: agents,
	}
}

// Run runs a workflow from a JSON definition.
func (r *PipelineRunner) Run(jsonWorkflow string) (map[string]string, error) {
	var workflow Workflow
	if err := json.Unmarshal([]byte(jsonWorkflow), &workflow); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow: %w", err)
	}

	outputs := make(map[string]string)
	previousOutput := ""

	for i, step := range workflow.Steps {
		agent, ok := r.agents[step.Agent]
		if !ok {
			return nil, fmt.Errorf("agent %q not found", step.Agent)
		}

		input := strings.ReplaceAll(step.Input, "{previous.output}", previousOutput)

		output, err := agent.Run(input)
		if err != nil {
			return nil, fmt.Errorf("step %d failed: %w", i+1, err)
		}

		stepName := fmt.Sprintf("step%d", i+1)
		outputs[stepName] = output
		previousOutput = output
	}

	return outputs, nil
}
