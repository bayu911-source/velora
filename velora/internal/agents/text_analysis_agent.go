
package agents

import (
	"velora/internal/services"
)

type TextAnalysisAgent struct {
	gemini *services.GeminiService
}

func NewTextAnalysisAgent(gemini *services.GeminiService) *TextAnalysisAgent {
	return &TextAnalysisAgent{
		gemini: gemini,
	}
}

func (a *TextAnalysisAgent) Name() string {
	return "text_analysis"
}

func (a *TextAnalysisAgent) Run(input string) (string, error) {
	prompt := "Analyze the following text: " + input
	return a.gemini.Generate(prompt, 0.5, 512)
}
