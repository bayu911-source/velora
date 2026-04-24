package agents

import (
	"context"
	"fmt"
)

// Run executes an agent and handles any errors that occur.
func Run(ctx context.Context, agent Agent, input string) (string, error) {
	output, err := agent.Execute(ctx, input)
	if err != nil {
		return "", fmt.Errorf("agent %q failed: %w", agent.Name(), err)
	}

	return output, nil
}
