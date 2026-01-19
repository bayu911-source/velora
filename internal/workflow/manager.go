
package workflow

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Manager is responsible for loading and managing workflows.
type Manager struct{}

// NewManager creates a new workflow manager.
func NewManager() *Manager {
	return &Manager{}
}

// Load loads a workflow from a YAML file.
func (m *Manager) Load(path string) (*Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var workflow Workflow
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return nil, err
	}

	return &workflow, nil
}
