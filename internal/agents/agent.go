
package agents

import (
	"context"
	"strings"
)

// Agent is the interface for all agents.
type Agent interface {
	Name() string
	Description() string
	Execute(ctx context.Context, input string) (string, error)
}

// SplitAgents splits a comma-separated string of agent names into a slice.
func SplitAgents(agentNames string) []string {
	return strings.Split(agentNames, ",")
}
