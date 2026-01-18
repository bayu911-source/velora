
package agents

// Agent is the interface for all AI agents.
type Agent interface {
	Name() string
	Run(input string) (string, error)
}
