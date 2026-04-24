package workflow

import (
	"context"
	"fmt"

	"velora/internal/agents"
)

// Engine is responsible for executing a workflow.
type Engine struct {
	agentRegistry *agents.Registry
	repo          *Repository
}

// NewEngine creates a new workflow engine.
func NewEngine(registry *agents.Registry, repo *Repository) *Engine {
	return &Engine{
		agentRegistry: registry,
		repo:          repo,
	}
}

// Run executes the given workflow.
func (e *Engine) Run(ctx context.Context, w *Workflow, initialInput string) (string, error) {
	w.State = StateRunning
	if err := e.repo.Save(w); err != nil {
		return "", fmt.Errorf("failed to save initial workflow state: %w", err)
	}

	var currentInput = initialInput
	var output string

	for _, step := range w.Steps() {
		agent, err := e.agentRegistry.Get(step.AgentName)
		if err != nil {
			return "", e.failWorkflow(w, fmt.Errorf("failed to get agent '%s': %w", step.AgentName, err))
		}

		fmt.Printf("Executing agent: %s\n", agent.Name())
		step.Input = currentInput

		output, err = agent.Execute(ctx, currentInput)
		if err != nil {
			step.Output = err.Error()
			_ = e.repo.Save(w) // Attempt to save error state
			return "", e.failWorkflow(w, fmt.Errorf("agent %s failed: %w", agent.Name(), err))
		}
		currentInput = output
		step.Output = output

		if err := e.repo.Save(w); err != nil {
			return "", e.failWorkflow(w, fmt.Errorf("failed to save workflow step: %w", err))
		}
	}

	w.State = StateCompleted
	if err := e.repo.Save(w); err != nil {
		return "", fmt.Errorf("failed to save final workflow state: %w", err)
	}

	return output, nil
}

func (e *Engine) failWorkflow(w *Workflow, err error) error {
	w.State = StateFailed
	if saveErr := e.repo.Save(w); saveErr != nil {
		return fmt.Errorf("error executing workflow: %v (and failed to save final state: %v)", err, saveErr)
	}
	return err
}
