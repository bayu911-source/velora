
package agents

import (
	"context"
	"fmt"

	"google.golang.org/genai"
	"velora/internal/services"
)

// TextAnalysisAgent performs analysis on a given text.
type TextAnalysisAgent struct {
	LLM services.LLM
}

// NewTextAnalysisAgent creates a new TextAnalysisAgent.
func NewTextAnalysisAgent(llm services.LLM) *TextAnalysisAgent {
	return &TextAnalysisAgent{
		LLM: llm,
	}
}

// Name returns the name of the agent.
func (a *TextAnalysisAgent) Name() string {
	return "text_analysis"
}

// Description returns a brief description of the agent.
func (a *TextAnalysisAgent) Description() string {
	return "Performs analysis on a given text, providing a summary, sentiment, and named entities."
}

// Run executes the agent's primary function: analyzing text.
func (a *TextAnalysisAgent) Run(ctx context.Context, input string) (string, error) {
	if a.LLM == nil {
		return "", fmt.Errorf("LLM service is not initialized")
	}

	// Create a specific prompt for the text analysis task.
	prompt := fmt.Sprintf("Analyze the following text and provide a summary of its key points, sentiment, and any named entities. Do not include any explanations, just the analysis.\n\nText: %s", input)

	// Generate the analysis using the LLM service.
	resp, err := a.LLM.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to analyze text: %w", err)
	}

	return resp, nil
}
