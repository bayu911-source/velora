
package agents

import (
	"context"
)

// CustomAgent is an agent created from a string prompt.
type CustomAgent struct {
	name        string
	prompt      string
	description string
	llm         LLMService
}

// NewCustomAgent creates a new CustomAgent.
func NewCustomAgent(name, description, prompt string, llm LLMService) *CustomAgent {
	return &CustomAgent{
		name:        name,
		description: description,
		prompt:      prompt,
		llm:         llm,
	}
}

// Name returns the name of the agent.
func (a *CustomAgent) Name() string {
	return a.name
}

// Description returns the description of the agent.
func (a *CustomAgent) Description() string {
	return a.description
}

// Execute executes the agent.
func (a *CustomAgent) Execute(ctx context.Context, input string) (string, error) {
	fullPrompt := a.prompt + "\n\n" + input
	return a.llm.Generate(ctx, fullPrompt)
}
