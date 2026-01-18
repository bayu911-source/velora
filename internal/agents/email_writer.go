
package agents

import (
	"velora/internal/services"
)

// EmailWriterAgent is an agent that writes emails.
type EmailWriterAgent struct {
	gemini *services.GeminiService
}

// NewEmailWriterAgent creates a new EmailWriterAgent.
func NewEmailWriterAgent(gemini *services.GeminiService) *EmailWriterAgent {
	return &EmailWriterAgent{
		gemini: gemini,
	}
}

// Name returns the name of the agent.
func (a *EmailWriterAgent) Name() string {
	return "email_writer"
}

// Run runs the agent.
func (a *EmailWriterAgent) Run(input string) (string, error) {
	prompt := "Write an email about: " + input
	return a.gemini.Generate(prompt, "gemini-2.5-pro", 0.7, 1024)
}
