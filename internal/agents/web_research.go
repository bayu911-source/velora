package agents

import (
	"context"
	"fmt"
)

// WebResearchAgent performs web research and summarizes the findings.
type WebResearchAgent struct {
	LLM LLMService
}

// NewWebResearchAgent creates a new WebResearchAgent.
func NewWebResearchAgent(llm LLMService) *WebResearchAgent {
	return &WebResearchAgent{
		LLM: llm,
	}
}

// Name returns the name of the agent.
func (a *WebResearchAgent) Name() string {
	return "web_research"
}

// Description returns a brief description of the agent.
func (a *WebResearchAgent) Description() string {
	return "Performs web research on a given topic and provides a detailed summary."
}

// Execute executes the agent's primary function: web research.
func (a *WebResearchAgent) Execute(ctx context.Context, input string) (string, error) {
	if a.LLM == nil {
		return "", fmt.Errorf("LLM service is not initialized")
	}

	// This is a simplified implementation. A real-world agent would use a search engine API.
	// For now, we'll simulate the research process by asking the LLM to act as a researcher.
	prompt := fmt.Sprintf("You are a world-class researcher. Research the following topic and provide a detailed summary of your findings. Include key facts, figures, and sources where possible. Do not include any explanations, just the research summary itself.\n\nTopic: %s", input)

	// Generate the research summary using the LLM service.
	resp, err := a.LLM.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to perform web research: %w", err)
	}

	return resp, nil
}
