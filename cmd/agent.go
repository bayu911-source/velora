
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"velora/config"
	"velora/internal/agents"
	"velora/internal/services"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage and run agents",
}

var runAgentCmd = &cobra.Command{
	Use:   "run [agent_name] [input]",
	Short: "Run a specific agent",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(".")
		if err != nil {
			log.Fatalf("cannot load config: %v", err)
		}

		agentName := args[0]
		input := args[1]

		// Initialize Gemini Service
		gemini, err := services.NewGeminiService(cfg.GeminiAPIKey)
		if err != nil {
			log.Fatalf("Failed to create Gemini service: %v", err)
		}

		// Register and get the agent
		agentRegistry := make(map[string]agents.Agent)
		agentRegistry["text_analysis"] = agents.NewTextAnalysisAgent(gemini)
		agentRegistry["email_writer"] = agents.NewEmailWriterAgent(gemini)
		agentRegistry["data_extractor"] = agents.NewDataExtractorAgent(gemini)
		agentRegistry["web_research"] = agents.NewWebResearchAgent(gemini)

		agent, ok := agentRegistry[agentName]
		if !ok {
			log.Fatalf("Agent %q not found", agentName)
		}

		output, err := agent.Run(input)
		if err != nil {
			log.Fatalf("Agent failed: %v", err)
		}

		fmt.Println(output)
	},
}

func init() {
	agentCmd.AddCommand(runAgentCmd)
	rootCmd.AddCommand(agentCmd)
}
