
package agents

import (
	"context"
)

// Run runs an agent.
func Run(ctx context.Context, agent Agent, input string) (string, error) {
	// In a real application, you might want to pass a services registry to the agent.
	return agent.Run(ctx, nil, input)
}
