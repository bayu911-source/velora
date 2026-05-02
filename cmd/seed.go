package cmd

import (
    "log"
    "time"

    "github.com/spf13/cobra"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"

    "velora/config"
    "velora/internal/database"
    "velora/internal/models"
    "velora/internal/repositories"
)

func NewSeedCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "seed",
        Short: "Seed demo data for Velora",
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

            tenantID := uuid.New()
            tenant := &models.Tenant{ID: tenantID, Name: "Velora Demo", Slug: "velora-demo", Plan: "pro", Status: "active", Branding: []byte(`{"theme":"velora"}`)}
            if err := repo.CreateTenant(tenant); err != nil {
                log.Fatalf("failed to create tenant: %v", err)
            }
            passwordHash, _ := bcrypt.GenerateFromPassword([]byte("demoPassword123"), bcrypt.DefaultCost)
            user := &models.User{BaseModel: models.BaseModel{ID: uuid.New(), TenantID: tenantID}, Email: "owner@velora.demo", Name: "Demo Owner", Role: "owner", PasswordHash: string(passwordHash), IsActive: true}
            if err := repo.CreateUser(user); err != nil {
                log.Fatalf("failed to create demo user: %v", err)
            }
            agent := &models.Agent{BaseModel: models.BaseModel{ID: uuid.New(), TenantID: tenantID}, Name: "Lead Generator", Type: "lead_generator", Description: "Generates warm leads using AI.", Prompt: "Generate a high-quality lead list for B2B SaaS companies.", Active: true}
            if err := repo.CreateAgent(agent); err != nil {
                log.Fatalf("failed to create demo agent: %v", err)
            }
            workflow := &models.Workflow{BaseModel: models.BaseModel{ID: uuid.New(), TenantID: tenantID}, Name: "New Lead Workflow", Description: "Demonstration workflow to generate leads and follow up.", Trigger: "lead_created", Actions: []byte(`[ {"name":"Lead Generator","agent_id":"` + agent.ID.String() + `","prompt":"Build a lead profile from the new lead input."} ]`), Status: "active"}
            if err := repo.CreateWorkflow(workflow); err != nil {
                log.Fatalf("failed to create demo workflow: %v", err)
            }
            log.Printf("demo tenant %s seeded successfully", tenant.ID.String())
            log.Printf("owner login: owner@velora.demo / demoPassword123")
            log.Printf("app url: http://localhost:%s", cfg.Port)
            time.Sleep(100 * time.Millisecond)
        },
    }
}
