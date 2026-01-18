
package cmd

import (
	"fmt"
	"os"

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
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
