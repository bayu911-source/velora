
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

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all saved workflows",
		RunE: func(cmd *cobra.Command, args []string) error {
			workflows, err := repo.ListAll()
			if err != nil {
				return fmt.Errorf("failed to list workflows: %w", err)
			}
			if len(workflows) == 0 {
				fmt.Println("No workflows found.")
				return nil
			}
			fmt.Println("Workflows:")
			for _, w := range workflows {
				fmt.Printf("- %s: %s (state=%s)\n", w.ID, w.Name, w.State)
			}
			return nil
		},
	}

	showCmd := &cobra.Command{
		Use:   "show [workflow-id]",
		Short: "Show details for a saved workflow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			w, err := repo.FindByID(args[0])
			if err != nil {
				return fmt.Errorf("failed to find workflow: %w", err)
			}
			if w == nil {
				return fmt.Errorf("workflow '%s' not found", args[0])
			}
			fmt.Printf("ID: %s\nName: %s\nState: %s\nCreated: %s\nUpdated: %s\n",
				w.ID, w.Name, w.State, w.CreatedAt.Format("2006-01-02 15:04:05"), w.UpdatedAt.Format("2006-01-02 15:04:05"))
			if len(w.Steps()) == 0 {
				fmt.Println("No workflow steps recorded.")
				return nil
			}
			fmt.Println("Steps:")
			for i, step := range w.Steps() {
				fmt.Printf("  %d. %s\n    Input: %s\n    Output: %s\n", i+1, step.AgentName, step.Input, step.Output)
			}
			return nil
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete [workflow-id]",
		Short: "Delete a workflow and its history",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := repo.DeleteByID(args[0]); err != nil {
				return fmt.Errorf("failed to delete workflow: %w", err)
			}
			fmt.Printf("Workflow '%s' deleted.\n", args[0])
			return nil
		},
	}

	workflowCmd.AddCommand(runCmd, createCmd, loadCmd, listCmd, showCmd, deleteCmd)

	return workflowCmd
}
