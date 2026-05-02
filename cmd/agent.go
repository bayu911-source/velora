
package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

func NewAgentCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "agent",
        Short: "Legacy agent CLI command stub",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Velora agent CLI is deprecated. Use the HTTP API for agent management.")
        },
    }
}
