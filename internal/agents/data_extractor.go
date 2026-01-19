
package agents

import (
	"context"
	"fmt"

	"velora/internal/services"
)

// DataExtractorAgent extracts structured data from unstructured text.
type DataExtractorAgent struct {
	llm services.LLMService
}

// NewDataExtractorAgent creates a new DataExtractorAgent.
func NewDataExtractorAgent(llm services.LLMService) *DataExtractorAgent {
	return &DataExtractorAgent{
		llm: llm,
	}
}

// Name returns the name of the agent.
func (a *DataExtractorAgent) Name() string {
	return "data_extractor"
}

// Description returns a brief description of the agent.
func (a *DataExtractorAgent) Description() string {
	return "Extracts structured data from unstructured text and formats it as JSON."
}

// Run executes the agent's primary function: extracting data.
func (a *DataExtractorAgent) Run(ctx context.Context, input string) (string, error) {
	if a.llm == nil {
		return "", fmt.Errorf("LLM service is not initialized")
	}

	// Create a specific prompt for the data extraction task.
	prompt := fmt.Sprintf("Extract the key information from the following text and format it as a JSON object. The keys of the JSON should be descriptive of the data being extracted. Do not include any explanations, just the JSON object itself.\n\nText: %s", input)

	// Generate the data using the LLM service.
	resp, err := a.llm.Generate(prompt, "gemini-1.5-pro", 0.7, 2048)
	if err != nil {
		return "", fmt.Errorf("failed to extract data: %w", err)
	}

	return resp, nil
}
