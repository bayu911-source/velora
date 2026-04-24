package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"velora/internal/agents"
	"velora/internal/server"
)

// NewServerCmd creates a new server command.
func NewServerCmd(registry *agents.Registry) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start the Velora HTTP server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Starting server on :8080")
			srv := server.NewServer(registry)
			if err := srv.ListenAndServe(); err != nil {
				log.Fatalf("could not start server: %v", err)
			}
		},
	}
}