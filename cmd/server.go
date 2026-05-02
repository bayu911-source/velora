package cmd

import (
    "fmt"
    "log"

    "github.com/spf13/cobra"
    "go.uber.org/zap"

    "velora/config"
    "velora/internal/database"
    "velora/internal/repositories"
    serverpkg "velora/internal/server"
)

func NewServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start the Velora HTTP server",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadConfig(".")
			if err != nil {
				log.Fatalf("failed to load config: %v", err)
			}

			db, err := database.Connect(cfg)
			if err != nil {
				log.Fatalf("failed to connect to database: %v", err)
			}

			repo := repositories.NewRepository(db)
			logger, _ := zap.NewProduction()
			app := serverpkg.NewApp(cfg, repo, logger)
			addr := fmt.Sprintf(":%s", cfg.Port)
			log.Printf("starting Velora API on %s", addr)
			if err := app.Listen(addr); err != nil {
				log.Fatalf("server exited: %v", err)
			}
		},
	}
}