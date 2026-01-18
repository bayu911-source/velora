
package pkg

import (
	// No imports needed for this file after consolidation
)

type WebResearchAgent struct {
	gemini *GeminiService
}

func NewWebResearchAgent(gemini *GeminiService) *WebResearchAgent {
	return &WebResearchAgent{
		gemini: gemini,
	}
}

func (a *WebResearchAgent) Name() string {
	return "web_research"
}

func (a *WebResearchAgent) Run(input string) (string, error) {
	// In a real implementation, this would use a search engine API
	prompt := "Research the web for: " + input
	return a.gemini.Generate(prompt, 0.2, 2048)
}
