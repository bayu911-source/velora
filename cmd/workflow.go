
package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"velora/internal/agents"
	"velora/internal/workflow"
)

// NewWorkflowCmd creates a new workflow command.
func NewWorkflowCmd(registry *agents.Registry) *cobra.Command {
	workflowCmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage and run workflows",
	}

	runCmd := &cobra.Command{
		Use:   "run [workflow-file]",
		Short: "Run a workflow from a file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			wm := workflow.NewManager()
			w, err := wm.Load(args[0])
			if err != nil {
				return err
			}

			we := workflow.NewEngine(registry)
			return we.Run(context.Background(), w)
		},
	}

	workflowCmd.AddCommand(runCmd)

	return workflowCmd
}
