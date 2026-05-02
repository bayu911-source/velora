
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

func NewWorkflowCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "workflow",
        Short: "Legacy workflow CLI command stub",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Velora workflow CLI is deprecated. Use the HTTP API for workflow management.")
        },
    }
}
