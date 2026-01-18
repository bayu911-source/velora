
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"velora/internal/agents"
)

func NewAgentCmd(registry *agents.Registry) *cobra.Command {
	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage and run agents",
	}

	runAgentCmd := &cobra.Command{
		Use:   "run [agent_name] [input]",
		Short: "Run a specific agent",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			agentName := args[0]
			var input string
			if len(args) > 1 {
				input = args[1]
			}

			agent, err := registry.Get(agentName)
			if err != nil {
				log.Fatalf("Agent %q not found", agentName)
			}

			output, err := agents.Run(cmd.Context(), agent, input)
			if err != nil {
				log.Fatalf("Agent failed: %v", err)
			}

			fmt.Println(output)
		},
	}

	listAgentsCmd := &cobra.Command{
		Use:   "list",
		Short: "List all available agents",
		Run: func(cmd *cobra.Command, args []string) {
			for _, agent := range registry.List() {
				fmt.Printf("- %s: %s\n", agent.Name(), agent.Description())
			}
		},
	}

	agentCmd.AddCommand(runAgentCmd)
	agentCmd.AddCommand(listAgentsCmd)

	return agentCmd
}
