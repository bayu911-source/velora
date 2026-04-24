package workflows

// Workflow represents a workflow with name and description.
type Workflow struct {
	Name        string
	Description string
}

// List returns a list of available workflows.
// For now, this is hardcoded. In the future, this could load from files or database.
func List() []Workflow {
	return []Workflow{
		{
			Name:        "Sample Code Generation Workflow",
			Description: "A simple workflow to generate Go code and then chat about it.",
		},
	}
}