package agents

import (
	"context"
	"fmt"
)

// EmailWriterAgent generates email copy.
type EmailWriterAgent struct {
	LLM LLMService
}

// NewEmailWriterAgent creates a new EmailWriterAgent.
func NewEmailWriterAgent(llm LLMService) *EmailWriterAgent {
	return &EmailWriterAgent{
		LLM: llm,
	}
}

// Name returns the name of the agent.
func (a *EmailWriterAgent) Name() string {
	return "email_writer"
}

// Description returns a brief description of the agent.
func (a *EmailWriterAgent) Description() string {
	return "Generates compelling email copy based on a given prompt."
}

// Execute executes the agent's primary function: writing email copy.
func (a *EmailWriterAgent) Execute(ctx context.Context, input string) (string, error) {
	if a.LLM == nil {
		return "", fmt.Errorf("LLM service is not initialized")
	}

	// Create a specific prompt for the email writing task.
	prompt := fmt.Sprintf("You are an expert copywriter. Write a compelling email based on the following prompt. The email should be clear, concise, and persuasive. Do not include any explanations, just the email copy itself.\n\nPrompt: %s", input)

	// Generate the email copy using the LLM service.
	resp, err := a.LLM.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to write email: %w", err)
	}

	return resp, nil
}
