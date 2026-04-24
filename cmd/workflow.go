
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"velora/internal/agents"
	"velora/internal/workflow"
	"velora/persistence"
)

// NewWorkflowCmd creates a new workflow command.
func NewWorkflowCmd(registry *agents.Registry) *cobra.Command {
	workflowCmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage and run workflows",
	}

	// Initialize repository and engine
	repo := workflow.NewRepository(persistence.DB)
	engine := workflow.NewEngine(registry, repo)

	runCmd := &cobra.Command{
		Use:   "run [workflow-id] [initial-input]",
		Short: "Run a workflow by its ID",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			w, err := repo.FindByID(args[0])
			if err != nil {
				return fmt.Errorf("failed to find workflow: %w", err)
			}
			if w == nil {
				return fmt.Errorf("workflow with ID '%s' not found", args[0])
			}

			output, err := engine.Run(context.Background(), w, args[1])
			if err != nil {
				return fmt.Errorf("workflow execution failed: %w", err)
			}

			fmt.Printf("Workflow completed successfully. Final output:\n%s\n", output)
			return nil
		},
	}

	createCmd := &cobra.Command{
		Use:   "create [name] [agent1,agent2,...]",
		Short: "Create a new workflow",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			agentNames := args[1]
			w := workflow.New(name, "", agents.SplitAgents(agentNames))
			if err := repo.Save(w); err != nil {
				return fmt.Errorf("failed to create workflow: %w", err)
			}
			fmt.Printf("Workflow '%s' created with ID: %s\n", w.Name, w.ID)
			return nil
		},
	}

	loadCmd := &cobra.Command{
		Use:   "load [file]",
		Short: "Load a workflow from a YAML file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := workflow.NewManager()
			w, err := manager.Load(args[0])
			if err != nil {
				return fmt.Errorf("failed to load workflow: %w", err)
			}
			if err := repo.Save(w); err != nil {
				return fmt.Errorf("failed to save workflow: %w", err)
			}
			fmt.Printf("Workflow '%s' loaded and saved with ID: %s\n", w.Name, w.ID)
			return nil
		},
	}

	workflowCmd.AddCommand(runCmd, createCmd, loadCmd)

	return workflowCmd
}
