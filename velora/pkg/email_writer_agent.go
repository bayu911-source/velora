
package pkg

import (
	// No imports needed for this file after consolidation
)

type EmailWriterAgent struct {
	gemini *GeminiService
}

func NewEmailWriterAgent(gemini *GeminiService) *EmailWriterAgent {
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
