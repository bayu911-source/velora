
package workflow

import (
	"fmt"
	"github.com/velora-chat/velora/pkg"
)

type Runner struct {
	agents map[string]pkg.Agent
}

func NewRunner() *Runner {
	return &Runner{
		agents: make(map[string]pkg.Agent),
	}
}

func (r *Runner) RegisterAgent(agent pkg.Agent) {
	r.agents[agent.Name()] = agent
}

func (r *Runner) Run(pipeline *Pipeline) (string, error) {
	var result string
	var err error

	for _, step := range pipeline.Steps {
		agent, ok := r.agents[step.Agent]
		if !ok {
			return "", fmt.Errorf("agent '%s' not found", step.Agent)
		}

		input := step.Input
		if result != "" {
			input = result // Use the output of the previous step as input
		}

		result, err = agent.Run(input)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}
