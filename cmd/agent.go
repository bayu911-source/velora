
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"velora/internal/agents"
	"velora/internal/plugin"
)

// NewAgentCmd creates a new agent command.
func NewAgentCmd(registry *agents.Registry, pluginManager *plugin.Manager) *cobra.Command {
	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage and run agents",
	}

	listAgentsCmd := &cobra.Command{
		Use:   "list",
		Short: "List all available agents",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Built-in Agents:")
			for _, agent := range registry.List() {
				fmt.Printf("- %s: %s\n", agent.Name(), agent.Description())
			}

			fmt.Println("\nPlugin Agents:")
			for _, agent := range pluginManager.Agents() {
				fmt.Printf("- %s\n", agent.Name())
			}
		},
	}

	agentCmd.AddCommand(listAgentsCmd)

	return agentCmd
}

// NewRunnableAgentCmd creates a new runnable agent command for a plugin.
func NewRunnableAgentCmd(agent plugin.Agent) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s [input]", agent.Name()),
		Short: fmt.Sprintf("Run the %s agent", agent.Name()),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Note: We're passing a background context here.
			// For more complex scenarios, you might want to handle context cancellation.
			output, err := agent.Execute(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("agent %s failed: %w", agent.Name(), err)
			}

			fmt.Println(output)
			return nil
		},
	}
}
