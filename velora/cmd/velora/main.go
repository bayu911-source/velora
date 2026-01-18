
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"velora/internal/agents"
	"velora/internal/config"
	"velora/internal/services"
	"velora/internal/utils"
	"velora/internal/workflow"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	pipelinePath := flag.String("pipeline", "", "Path to the pipeline file")
	logPath := flag.String("log", "velora.log", "Path to the log file")
	flag.Parse()

	logger, err := utils.NewLogger(*logPath)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	geminiService, err := services.NewGeminiService(cfg.GeminiModel)
	if err != nil {
		logger.Fatalf("Failed to create Gemini service: %v", err)
	}

	textAnalysisAgent := agents.NewTextAnalysisAgent(geminiService)
	emailWriterAgent := agents.NewEmailWriterAgent(geminiService)
	webResearchAgent := agents.NewWebResearchAgent(geminiService)

	runner := workflow.NewRunner()
	runner.RegisterAgent(textAnalysisAgent)
	runner.RegisterAgent(emailWriterAgent)
	runner.RegisterAgent(webResearchAgent)

	if *pipelinePath != "" {
		pipeline, err := workflow.LoadPipeline(*pipelinePath)
		if err != nil {
			logger.Fatalf("Failed to load pipeline: %v", err)
		}

		// For now, let's assume a simple, non-scheduled run
		result, err := runner.Run(pipeline)
		if err != nil {
			logger.Fatalf("Failed to run pipeline: %v", err)
		}
		fmt.Println(result)
	} else {
		// Example of using the scheduler
		scheduler := workflow.NewScheduler()
		pipeline := &workflow.Pipeline{
			Name: "Example Scheduled Pipeline",
			Steps: []workflow.Step{
				{Agent: "text_analysis", Input: "Analyze the latest AI trends."},
			},
		}
		scheduler.Schedule("ai_trends_analysis", pipeline, 1*time.Hour)
		scheduler.Start("ai_trends_analysis", runner)

		logger.Println("Scheduler started. Press Ctrl+C to exit.")
		// Keep the application running
		select {}
	}
}
