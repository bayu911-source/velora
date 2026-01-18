
package agents

import (
	"velora/internal/services"
)

// CustomAgent is an agent created from a string prompt.
type CustomAgent struct {
	name   string
	prompt string
	gemini *services.GeminiService
}

// NewCustomAgent creates a new CustomAgent.
func NewCustomAgent(name, prompt string, gemini *services.GeminiService) *CustomAgent {
	return &CustomAgent{
		name:   name,
		prompt: prompt,
		gemini: gemini,
	}
}

// Name returns the name of the agent.
func (a *CustomAgent) Name() string {
	return a.name
}

// Run runs the agent.
func (a *CustomAgent) Run(input string) (string, error) {
	fullPrompt := a.prompt + "\n\n" + input
	return a.gemini.Generate(fullPrompt, "gemini-2.5-pro", 0.7, 1024)
}
