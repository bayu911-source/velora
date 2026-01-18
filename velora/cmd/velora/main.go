
package main

import (
	"fmt"
	"log"

	"github.com/velora-chat/velora/internal/workflow"
	"github.com/velora-chat/velora/pkg"
)

func main() {
	geminiService, err := pkg.NewGeminiService("gemini-1.5-pro-latest")
	if err != nil {
		log.Fatalf("Failed to create Gemini service: %v", err)
	}

	// Create agents
	webResearchAgent := pkg.NewWebResearchAgent(geminiService)
	textAnalysisAgent := pkg.NewTextAnalysisAgent(geminiService)
	emailWriterAgent := pkg.NewEmailWriterAgent(geminiService)

	// Create a new runner and register the agents
	runner := workflow.NewRunner()
	runner.RegisterAgent(webResearchAgent)
	runner.RegisterAgent(textAnalysisAgent)
	runner.RegisterAgent(emailWriterAgent)

	// Load the pipeline
	pipeline, err := workflow.LoadPipeline("email_pipeline.json")
	if err != nil {
		log.Fatalf("Failed to load pipeline: %v", err)
	}

	// Run the pipeline
	result, err := runner.Run(pipeline)
	if err != nil {
		log.Fatalf("Failed to run pipeline: %v", err)
	}

	// Print the result
	fmt.Println(result)
}
