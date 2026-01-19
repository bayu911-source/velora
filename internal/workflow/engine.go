
package workflow

import (
	"context"
	"fmt"

	"velora/internal/agents"
)

// Engine is responsible for executing workflows.
type Engine struct {
	registry *agents.Registry
}

// NewEngine creates a new workflow engine.
func NewEngine(registry *agents.Registry) *Engine {
	return &Engine{
		registry: registry,
	}
}

// Run executes the given workflow.
func (e *Engine) Run(ctx context.Context, workflow *Workflow) error {
	for _, step := range workflow.Steps {
		agent, err := e.registry.Get(step.Agent)
		if err != nil {
			return fmt.Errorf("agent '%s' not found: %w", step.Agent, err)
		}

		fmt.Printf("Running step: %s\n", step.Name)
		output, err := agent.Run(ctx, e.registry.LLM(), step.Input)
		if err != nil {
			return fmt.Errorf("step '%s' failed: %w", step.Name, err)
		}

		fmt.Printf("Step '%s' output: %s\n", step.Name, output)
	}

	return nil
}
