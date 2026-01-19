
package agents

import (
	"context"

	"velora/internal/services"
)

// Agent is the interface for all agents.
type Agent interface {
	Name() string
	Description() string
	Run(ctx context.Context, input string) (string, error)
}
