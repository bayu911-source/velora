
package workflow

// Workflow defines a series of steps to be executed.
type Workflow struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Steps       []Step `yaml:"steps"`
}

// Step defines a single step in a workflow.
type Step struct {
	Name  string `yaml:"name"`
	Agent string `yaml:"agent"`
	Input string `yaml:"input"`
}
