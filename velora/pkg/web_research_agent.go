
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

func (a *WebResearchAgent) Run(memory *MemoryManager, input string) (string, error) {
	// In a real implementation, this would use a search engine API
	// For now, we'll just pass the input to the Gemini service.
	// We could store the result in memory if needed:
	// memory.Set("research_result", result)
	prompt := "Research the web for: " + input
	return a.gemini.Generate(prompt, 0.2, 2048)
}
