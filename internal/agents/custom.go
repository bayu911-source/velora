
package agents

import (
	"context"
	"velora/internal/services"
)

// CustomAgent is an agent created from a string prompt.
type CustomAgent struct {
	name        string
	prompt      string
	description string
	llm         services.LLMService
}

// NewCustomAgent creates a new CustomAgent.
func NewCustomAgent(name, description, prompt string, llm services.LLMService) *CustomAgent {
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

// Run runs the agent.
func (a *CustomAgent) Run(ctx context.Context, input string) (string, error) {
	fullPrompt := a.prompt + "\n\n" + input
	return a.llm.Generate(fullPrompt, "gemini-1.5-pro", 0.7, 1024)
}
