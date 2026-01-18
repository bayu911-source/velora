
package agents

import "context"

// Agent is the interface for all AI agents.
type Agent interface {
	Name() string
	Description() string
	Run(ctx context.Context, input string) (string, error)
}
