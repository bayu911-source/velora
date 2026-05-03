
package agents

import (
	"context"
	"strings"
)

// LLMService defines the interface for LLM operations used by agents.
type LLMService interface {
	Generate(ctx context.Context, prompt string) (string, error)
	Close()
}

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
