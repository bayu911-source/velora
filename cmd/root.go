package cmd

import (
	"fmt"
	"log"
	"os"

	"velora/config"
	"velora/internal/agents"
	"velora/internal/plugin"
	"velora/internal/services"
	"velora/persistence"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "velora",
	Short: "Velora is a framework for building and managing AI agents.",
	Long:  `Velora is a flexible and extensible framework for building and managing AI agents and workflows.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action when no subcommand is provided
		cmd.Help()
	},
}

// Execute executes the root command.
func Execute(pluginManager *plugin.Manager) {
	// Initialize the database
	persistence.InitDB()

	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize LLM service
	llm, err := services.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize LLM service: %v", err)
	}
	defer llm.Close()

	// Initialize agent registry
	registry := agents.NewRegistry(llm)

	// Register built-in agents
	registry.Register(agents.NewCodeGenerator(llm))
	registry.Register(agents.NewChatAgent(llm))

	// Add agent command
	agentCmd := NewAgentCmd(registry, pluginManager)

	// Register agents from plugins as sub-commands
	for _, agent := range pluginManager.Agents() {
		agentCmd.AddCommand(NewRunnableAgentCmd(agent))
	}

	rootCmd.AddCommand(agentCmd)

	// Add workflow command
	rootCmd.AddCommand(NewWorkflowCmd(registry))

	// Add server command
	rootCmd.AddCommand(NewServerCmd(registry))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
