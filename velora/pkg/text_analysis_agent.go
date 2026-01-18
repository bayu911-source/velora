
package pkg

import (
	// No imports needed for this file after consolidation
)

type TextAnalysisAgent struct {
	gemini *GeminiService
}

func NewTextAnalysisAgent(gemini *GeminiService) *TextAnalysisAgent {
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
