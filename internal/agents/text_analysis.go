
package agents

import (
	"velora/internal/services"
)

// TextAnalysisAgent is an agent that analyzes text.
type TextAnalysisAgent struct {
	gemini *services.GeminiService
}

// NewTextAnalysisAgent creates a new TextAnalysisAgent.
func NewTextAnalysisAgent(gemini *services.GeminiService) *TextAnalysisAgent {
	return &TextAnalysisAgent{
		gemini: gemini,
	}
}

// Name returns the name of the agent.
func (a *TextAnalysisAgent) Name() string {
	return "text_analysis"
}

// Run runs the agent.
func (a *TextAnalysisAgent) Run(input string) (string, error) {
	prompt := "Analyze the following text: " + input
	return a.gemini.Generate(prompt, "gemini-2.5-pro", 0.7, 1024)
}
