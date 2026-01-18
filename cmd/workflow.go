
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"velora/internal/agents"
	"velora/internal/workflows"
)

// NewWorkflowCmd membuat perintah cobra baru untuk mengelola alur kerja.
func NewWorkflowCmd(registry *agents.Registry) *cobra.Command {
	workflowCmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage and run workflows",
	}

	createWorkflowCmd := &cobra.Command{
		Use:   "create [name] [agent1] [agent2] ...",
		Short: "Create a new workflow",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			agentNames := args[1:]

			// Pastikan semua agen ada di registri.
			for _, agentName := range agentNames {
				if _, err := registry.Get(agentName); err != nil {
					log.Fatalf("Agen '%s' tidak ditemukan", agentName)
				}
			}

			w, err := workflows.New(name, agentNames)
			if err != nil {
				log.Fatalf("Gagal membuat alur kerja: %v", err)
			}

			fmt.Printf("Alur kerja '%s' dibuat dengan ID: %s\n", w.Name, w.ID)
		},
	}

	runWorkflowCmd := &cobra.Command{
		Use:   "run [workflow_id] [input]",
		Short: "Run a specific workflow",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			workflowID := args[0]
			input := args[1]

			w, err := workflows.GetWorkflowByID(workflowID)
			if err != nil {
				log.Fatalf("Gagal mengambil alur kerja: %v", err)
			}
			if w == nil {
				log.Fatalf("Alur kerja dengan ID '%s' tidak ditemukan", workflowID)
			}

			output, err := w.Run(cmd.Context(), registry, input)
			if err != nil {
				log.Fatalf("Alur kerja gagal: %v", err)
			}

			fmt.Println(output)
		},
	}

	listWorkflowsCmd := &cobra.Command{
		Use:   "list",
		Short: "List all workflows",
		Run: func(cmd *cobra.Command, args []string) {
			ws, err := workflows.ListWorkflows()
			if err != nil {
				log.Fatalf("Gagal mengambil daftar alur kerja: %v", err)
			}

			for _, w := range ws {
				fmt.Printf("- %s: %s (%s)\n", w.ID, w.Name, w.State)
			}
		},
	}

	workflowCmd.AddCommand(createWorkflowCmd)
	workflowCmd.AddCommand(runWorkflowCmd)
	workflowCmd.AddCommand(listWorkflowsCmd)

	return workflowCmd
}
