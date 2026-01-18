
package agents

import (
	"velora/internal/services"
)

// DataExtractorAgent is an agent that extracts data from text.
type DataExtractorAgent struct {
	gemini *services.GeminiService
}

// NewDataExtractorAgent creates a new DataExtractorAgent.
func NewDataExtractorAgent(gemini *services.GeminiService) *DataExtractorAgent {
	return &DataExtractorAgent{
		gemini: gemini,
	}
}

// Name returns the name of the agent.
func (a *DataExtractorAgent) Name() string {
	return "data_extractor"
}

// Run runs the agent.
func (a *DataExtractorAgent) Run(input string) (string, error) {
	prompt := "Extract data from the following text: " + input
	return a.gemini.Generate(prompt, "gemini-2.5-pro", 0.7, 1024)
}
