
package agents

import (
	"velora/internal/services"
)

type EmailWriterAgent struct {
	gemini *services.GeminiService
}

func NewEmailWriterAgent(gemini *services.GeminiService) *EmailWriterAgent {
	return &EmailWriterAgent{
		gemini: gemini,
	}
}

func (a *EmailWriterAgent) Name() string {
	return "email_writer"
}

func (a *EmailWriterAgent) Run(input string) (string, error) {
	prompt := "Write an email about: " + input
	return a.gemini.Generate(prompt, 0.7, 1024)
}
