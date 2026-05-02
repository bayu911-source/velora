package cmd

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/spf13/cobra"
    "go.uber.org/zap"

    "velora/config"
    "velora/internal/database"
    "velora/internal/repositories"
    "velora/internal/services"
    "velora/internal/workers"
)

func NewWorkerCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "worker",
        Short: "Start Velora background worker",
        Run: func(cmd *cobra.Command, args []string) {
            cfg, err := config.LoadConfig(".")
            if err != nil {
                log.Fatalf("failed to load config: %v", err)
            }
            db, err := database.Connect(cfg)
            if err != nil {
                log.Fatalf("failed to connect database: %v", err)
            }
            repo := repositories.NewRepository(db)
            logger, _ := zap.NewProduction()
            services := services.NewServices(cfg, repo, logger)
            worker := workers.NewWorkerServer(cfg.RedisURL, services.Workflow, logger)

            ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
            defer stop()
            if err := worker.Start(ctx); err != nil {
                log.Fatalf("worker failed: %v", err)
            }
        },
    }
}
