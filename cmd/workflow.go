
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"velora/config"
	"velora/internal/agents"
	"velora/internal/services"
	"velora/internal/workflow"
)

var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Manage and run workflows",
}

var runWorkflowCmd = &cobra.Command{
	Use:   "run [json_workflow]",
	Short: "Run a workflow from a JSON definition",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(".")
		if err != nil {
			log.Fatalf("cannot load config: %v", err)
		}

		jsonWorkflow := args[0]

		// Initialize Gemini Service
		gemini, err := services.NewGeminiService(cfg.GeminiAPIKey)
		if err != nil {
			log.Fatalf("Failed to create Gemini service: %v", err)
		}

		// Register Agents
		agentRegistry := make(map[string]agents.Agent)
		agentRegistry["text_analysis"] = agents.NewTextAnalysisAgent(gemini)
		agentRegistry["email_writer"] = agents.NewEmailWriterAgent(gemini)
		agentRegistry["data_extractor"] = agents.NewDataExtractorAgent(gemini)
		agentRegistry["web_research"] = agents.NewWebResearchAgent(gemini)

		// Initialize Pipeline Runner
		pipelineRunner := workflow.NewPipelineRunner(agentRegistry)

		outputs, err := pipelineRunner.Run(jsonWorkflow)
		if err != nil {
			log.Fatalf("Workflow failed: %v", err)
		}

		for step, output := range outputs {
			fmt.Printf("%s: %s\n", step, output)
		}
	},
}

func init() {
	workflowCmd.AddCommand(runWorkflowCmd)
	rootCmd.AddCommand(workflowCmd)
}
