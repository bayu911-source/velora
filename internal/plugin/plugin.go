
package plugin

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
)

// Agent is the interface that all agent plugins must implement.
// Agents are the basic building blocks of Velora. They are responsible for
// performing a specific task.
type Agent interface {
	// Name returns the name of the agent.
	Name() string
	// Execute executes the agent with the given input and returns the output.
	Execute(ctx context.Context, input string) (string, error)
}

// Manager is responsible for loading and managing agent plugins.
type Manager struct {
	agents map[string]Agent
}

// NewManager creates a new plugin manager.
func NewManager() *Manager {
	return &Manager{
		agents: make(map[string]Agent),
	}
}

// LoadAgentsFromDir loads all agent plugins from a directory.
// It looks for shared object files (.so) and loads them as plugins.
func (m *Manager) LoadAgentsFromDir(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".so" {
			continue
		}

		path := filepath.Join(dir, file.Name())
		plug, err := plugin.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open plugin %s: %w", path, err)
		}

		sym, err := plug.Lookup("Agent")
		if err != nil {
			return fmt.Errorf("failed to lookup Agent symbol in %s: %w", path, err)
		}

		agent, ok := sym.(Agent)
		if !ok {
			return fmt.Errorf("unexpected type from symbol in %s: expected Agent", path)
		}

		m.agents[agent.Name()] = agent
	}

	return nil
}

// Agents returns a list of all loaded agents.
func (m *Manager) Agents() []Agent {
	var agents []Agent
	for _, agent := range m.agents {
		agents = append(agents, agent)
	}
	return agents
}

// GetAgent returns an agent by name.
func (m *Manager) GetAgent(name string) (Agent, bool) {
	agent, ok := m.agents[name]
	return agent, ok
}
