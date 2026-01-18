
package agents

import (
	"velora/internal/services"
)

type WebResearchAgent struct {
	gemini *services.GeminiService
}

func NewWebResearchAgent(gemini *services.GeminiService) *WebResearchAgent {
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
