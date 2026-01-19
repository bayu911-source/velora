
package workflow

import (
	"encoding/json"
	"fmt"
	"strings"
	"velora/internal/agents"
	"velora/persistence"
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
	var workflowDef Workflow
	if err := json.Unmarshal([]byte(jsonWorkflow), &workflowDef); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow: %w", err)
	}

	// Create a new workflow in the database.
	wf, err := persistence.CreateWorkflow(workflowDef.Name, "running")
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow record: %w", err)
	}

	outputs := make(map[string]string)
	previousOutput := ""

	for i, step := range workflowDef.Steps {
		agent, ok := r.agents[step.Agent]
		if !ok {
			r.updateWorkflowState(wf.ID, "failed")
			return nil, fmt.Errorf("agent %q not found", step.Agent)
		}

		input := strings.ReplaceAll(step.Input, "{previous.output}", previousOutput)

		output, err := agent.Run(input)
		if err != nil {
			r.updateWorkflowState(wf.ID, "failed")
			persistence.CreateWorkflowStep(wf.ID, step.Agent, input, fmt.Sprintf("Error: %v", err))
			return nil, fmt.Errorf("step %d failed: %w", i+1, err)
		}

		// Record the step execution.
		if err := persistence.CreateWorkflowStep(wf.ID, step.Agent, input, output); err != nil {
			// Log the error but continue the workflow.
			fmt.Printf("warning: failed to record workflow step: %v\n", err)
		}

		stepName := fmt.Sprintf("step%d", i+1)
		outputs[stepName] = output
		previousOutput = output
	}

	// Update the workflow state to "completed".
	if err := r.updateWorkflowState(wf.ID, "completed"); err != nil {
		// Log the error but consider the workflow successful from the user's perspective.
		fmt.Printf("warning: failed to update workflow state to completed: %v\n", err)
	}

	return outputs, nil
}

func (r *PipelineRunner) updateWorkflowState(id, state string) {
	if err := persistence.UpdateWorkflowState(id, state); err != nil {
		fmt.Printf("warning: failed to update workflow state: %v\n", err)
	}
}
