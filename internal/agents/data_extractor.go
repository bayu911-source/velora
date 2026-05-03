
package agents

import (
	"context"
	"fmt"
)

// DataExtractorAgent extracts structured data from unstructured text.
type DataExtractorAgent struct {
	LLM LLMService
}

// NewDataExtractorAgent creates a new DataExtractorAgent.
func NewDataExtractorAgent(llm LLMService) *DataExtractorAgent {
	return &DataExtractorAgent{
		LLM: llm,
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

// Execute executes the agent's primary function: extracting data.
func (a *DataExtractorAgent) Execute(ctx context.Context, input string) (string, error) {
	if a.LLM == nil {
		return "", fmt.Errorf("LLM service is not initialized")
	}

	// Create a specific prompt for the data extraction task.
	prompt := fmt.Sprintf("Extract the key information from the following text and format it as a JSON object. The keys of the JSON should be descriptive of the data being extracted. Do not include any explanations, just the JSON object itself.\n\nText: %s", input)

	// Generate the data using the LLM service.
	resp, err := a.LLM.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to extract data: %w", err)
	}

	return resp, nil
}
