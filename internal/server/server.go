package server

import (
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/compress"
    "github.com/gofiber/fiber/v2/middleware/recover"
    "go.uber.org/zap"

    "velora/config"
    "velora/internal/handlers"
    "velora/internal/middleware"
    "velora/internal/repositories"
    "velora/internal/services"
)

// NewApp constructs a Fiber application with SaaS routes and middleware.
func NewApp(cfg *config.Config, dbRepo *repositories.Repository, logger *zap.Logger) *fiber.App {
    svc := services.NewServices(cfg, dbRepo, logger)
    app := fiber.New(fiber.Config{AppName: cfg.AppName, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second})
    app.Use(middleware.CORS())
    app.Use(recover.New())
    app.Use(compress.New())
    app.Use(middleware.RateLimiter(cfg.RateLimitMax, cfg.RateLimitWindow))

    h := handlers.NewHandler(svc, cfg, dbRepo, logger)
    api := app.Group("/api/v1")

    auth := api.Group("/auth")
    auth.Post("/register", h.Register)
    auth.Post("/login", h.Login)
    auth.Post("/refresh", h.Refresh)

    tenants := api.Group("/tenants")
    tenants.Post("", h.CreateTenant)
    tenants.Get("/:tenantId", middleware.Auth(cfg.JWTSecret), h.GetTenant)
    tenants.Post("/:tenantId/invite", middleware.Auth(cfg.JWTSecret), middleware.Role(services.RoleOwner, services.RoleAdmin), h.InviteUser)

    tenantGroup := api.Group("", middleware.Auth(cfg.JWTSecret), middleware.Tenant())
    agents := tenantGroup.Group("/agents")
    agents.Post("", h.CreateAgent)
    agents.Get("", h.ListAgents)
    agents.Get("/:id", h.GetAgent)
    agents.Put("/:id", h.UpdateAgent)
    agents.Delete("/:id", h.DeleteAgent)
    agents.Post("/:id/run", h.RunAgent)

    workflows := tenantGroup.Group("/workflows")
    workflows.Post("", h.CreateWorkflow)
    workflows.Get("", h.ListWorkflows)
    workflows.Get("/:id", h.GetWorkflow)
    workflows.Delete("/:id", h.DeleteWorkflow)
    workflows.Post("/:id/run", h.RunWorkflow)

    leads := tenantGroup.Group("/leads")
    leads.Post("", h.CreateLead)
    leads.Get("", h.ListLeads)
    leads.Get("/:id", h.GetLead)
    leads.Post("/:id/score", h.ScoreLead)

    integrations := tenantGroup.Group("/integrations")
    integrations.Post("", h.CreateIntegration)
    integrations.Get("", h.ListIntegrations)

    billing := tenantGroup.Group("/billing")
    billing.Get("/plans", h.GetPlans)
    billing.Post("/subscriptions", h.CreateSubscription)
    billing.Post("/invoices", h.CreateInvoice)

    admin := api.Group("/admin", middleware.Auth(cfg.JWTSecret), middleware.Role(services.RoleOwner))
    admin.Get("/tenants", h.ListTenants)
    admin.Post("/tenants/:id/suspend", h.SuspendTenant)
    admin.Get("/logs", h.ListAuditLogs)
    admin.Get("/analytics", h.GetAnalytics)

    app.Get("/health", h.Health)
    app.Get("/", func(c *fiber.Ctx) error { return c.SendString("Velora AI Automation Framework API is running") })

    return app
}
