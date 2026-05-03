package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "velora",
    Short: "Velora AI Automation Framework",
    Long:  `Velora is a multi-tenant AI automation SaaS backend built with Go, Fiber, GORM, Redis, and Asynq.`,
    Run: func(cmd *cobra.Command, args []string) {
        cmd.Help()
    },
}

func Execute() {
    rootCmd.AddCommand(NewServerCmd())
    rootCmd.AddCommand(NewWorkerCmd())
    rootCmd.AddCommand(NewSeedCmd())
    rootCmd.AddCommand(NewWorkflowCmd())

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
