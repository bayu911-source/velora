package handlers

import (
    "encoding/json"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/go-playground/validator/v10"
    "github.com/google/uuid"
    "github.com/hibiken/asynq"
    "go.uber.org/zap"
    "gorm.io/datatypes"

    "velora/config"
    "velora/internal/dto"
    "velora/internal/models"
    "velora/internal/pkg/utils/response"
    "velora/internal/repositories"
    "velora/internal/services"
)

// Handler routes service calls to HTTP endpoints.
type Handler struct {
    services      *services.Services
    cfg           *config.Config
    validator     *validator.Validate
    repo          *repositories.Repository
    scheduler     *asynq.Scheduler
    logger        *zap.Logger
}

func NewHandler(svc *services.Services, cfg *config.Config, repo *repositories.Repository, logger *zap.Logger) *Handler {
    return &Handler{
        services:  svc,
        cfg:       cfg,
        repo:      repo,
        validator: validator.New(),
        logger:    logger,
    }
}

func (h *Handler) validate(c *fiber.Ctx, body any) error {
    if err := c.BodyParser(body); err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    if err := h.validator.Struct(body); err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    return nil
}

func (h *Handler) Health(c *fiber.Ctx) error {
    return response.JSON(c, fiber.StatusOK, map[string]string{"status": "ok", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}

func (h *Handler) Register(c *fiber.Ctx) error {
    var payload dto.RegisterRequest
    if err := h.validate(c, &payload); err != nil {
        return err
    }
    tenant, user, accessToken, refreshToken, err := h.services.Auth.RegisterTenant(payload.TenantName, payload.Name, payload.Email, payload.Password)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    return response.JSON(c, fiber.StatusCreated, map[string]any{"tenant": tenant, "user": user, "access_token": accessToken, "refresh_token": refreshToken})
}

func (h *Handler) Login(c *fiber.Ctx) error {
    var payload dto.LoginRequest
    if err := h.validate(c, &payload); err != nil {
        return err
    }
    tenantID, err := uuid.Parse(payload.TenantID)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    user, accessToken, refreshToken, err := h.services.Auth.Login(tenantID, payload.Email, payload.Password)
    if err != nil {
        return response.Error(c, fiber.StatusUnauthorized, err)
    }
    return response.JSON(c, fiber.StatusOK, map[string]any{"user": user, "access_token": accessToken, "refresh_token": refreshToken})
}

func (h *Handler) Refresh(c *fiber.Ctx) error {
    var payload dto.RefreshRequest
    if err := h.validate(c, &payload); err != nil {
        return err
    }
    accessToken, refreshToken, err := h.services.Auth.RefreshTokens(payload.RefreshToken)
    if err != nil {
        return response.Error(c, fiber.StatusUnauthorized, err)
    }
    return response.JSON(c, fiber.StatusOK, map[string]any{"access_token": accessToken, "refresh_token": refreshToken})
}

func (h *Handler) CreateTenant(c *fiber.Ctx) error {
    var payload dto.TenantRequest
    if err := h.validate(c, &payload); err != nil {
        return err
    }
    tenant, err := h.services.Tenant.CreateTenant(payload.Name)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    return response.JSON(c, fiber.StatusCreated, tenant)
}

func (h *Handler) GetTenant(c *fiber.Ctx) error {
    id := c.Params("tenantId")
    tenantID, err := uuid.Parse(id)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    tenant, err := h.services.Tenant.GetTenant(tenantID)
    if err != nil {
        return response.Error(c, fiber.StatusNotFound, err)
    }
    return response.JSON(c, fiber.StatusOK, tenant)
}

func (h *Handler) InviteUser(c *fiber.Ctx) error {
    var payload dto.InviteUserRequest
    if err := h.validate(c, &payload); err != nil {
        return err
    }
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    user, err := h.services.Auth.InviteUser(tenantUUID, payload.Email, payload.Name, payload.Role)
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusCreated, user)
}

func (h *Handler) CreateAgent(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    var payload dto.AgentRequest
    if err := h.validate(c, &payload); err != nil {
        return err
    }
    agent := &models.Agent{
        BaseModel:   models.BaseModel{ID: uuid.New(), TenantID: tenantUUID},
        Name:        payload.Name,
        Type:        payload.Type,
        Description: payload.Description,
        Prompt:      payload.Prompt,
        Active:      true,
    }
    if err := h.services.Agent.CreateAgent(agent); err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    return response.JSON(c, fiber.StatusCreated, agent)
}

func (h *Handler) ListAgents(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    agents, err := h.services.Agent.ListAgents(tenantUUID)
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusOK, agents)
}

func (h *Handler) GetAgent(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    agentID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    agent, err := h.services.Agent.GetAgent(tenantUUID, agentID)
    if err != nil {
        return response.Error(c, fiber.StatusNotFound, err)
    }
    return response.JSON(c, fiber.StatusOK, agent)
}

func (h *Handler) UpdateAgent(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    agentID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    var payload dto.AgentRequest
    if err := h.validate(c, &payload); err != nil {
        return err
    }
    agent, err := h.services.Agent.GetAgent(tenantUUID, agentID)
    if err != nil {
        return response.Error(c, fiber.StatusNotFound, err)
    }
    agent.Name = payload.Name
    agent.Type = payload.Type
    agent.Description = payload.Description
    agent.Prompt = payload.Prompt
    if err := h.services.Agent.UpdateAgent(agent); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusOK, agent)
}

func (h *Handler) DeleteAgent(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    agentID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    if err := h.services.Agent.DeleteAgent(tenantUUID, agentID); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.Message(c, fiber.StatusNoContent, "agent deleted")
}

func (h *Handler) RunAgent(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    agentID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    var payload dto.AgentRunRequest
    if err := h.validate(c, &payload); err != nil {
        return err
    }
    output, err := h.services.Agent.RunAgent(tenantUUID, agentID, payload.Input)
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusOK, map[string]any{"output": output})
}

func (h *Handler) CreateWorkflow(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    var payload dto.WorkflowRequest
    if err := h.validate(c, &payload); err != nil {
        return err
    }
    actionsBytes, err := json.Marshal(payload.Actions)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    workflow := &models.Workflow{
        BaseModel:   models.BaseModel{ID: uuid.New(), TenantID: tenantUUID},
        Name:        payload.Name,
        Description: payload.Description,
        Trigger:     payload.Trigger,
        Actions:     actionsBytes,
        Schedule:    payload.Schedule,
        Status:      "active",
    }
    if err := h.services.Workflow.CreateWorkflow(workflow); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    if payload.Schedule != "" {
        _ = h.services.Workflow.ScheduleWorkflow(workflow)
    }
    return response.JSON(c, fiber.StatusCreated, workflow)
}

func (h *Handler) ListWorkflows(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    workflows, err := h.services.Workflow.ListWorkflows(tenantUUID)
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusOK, workflows)
}

func (h *Handler) GetWorkflow(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    workflowID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    workflow, err := h.services.Workflow.GetWorkflow(tenantUUID, workflowID)
    if err != nil {
        return response.Error(c, fiber.StatusNotFound, err)
    }
    return response.JSON(c, fiber.StatusOK, workflow)
}

func (h *Handler) DeleteWorkflow(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    workflowID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    if err := h.services.Workflow.DeleteWorkflow(tenantUUID, workflowID); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.Message(c, fiber.StatusNoContent, "workflow deleted")
}

func (h *Handler) RunWorkflow(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    workflowID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    var payload dto.WorkflowRunRequest
    if err := h.validate(c, &payload); err != nil {
        return err
    }
    if err := h.services.Workflow.EnqueueWorkflowRun(tenantUUID, workflowID, payload.Input); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusAccepted, map[string]string{"status": "queued"})
}

func (h *Handler) CreateLead(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    var payload dto.LeadRequest
    if err := h.validate(c, &payload); err != nil {
        return err
    }
    lead := &models.Lead{
        BaseModel: models.BaseModel{ID: uuid.New(), TenantID: tenantUUID},
        Name:      payload.Name,
        Email:     payload.Email,
        Company:   payload.Company,
        Status:    "new",
    }
    if err := h.services.CRM.CreateLead(lead); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusCreated, lead)
}

func (h *Handler) ListLeads(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    leads, err := h.services.CRM.ListLeads(tenantUUID)
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusOK, leads)
}

func (h *Handler) GetLead(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    leadID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    lead, err := h.services.CRM.GetLead(tenantUUID, leadID)
    if err != nil {
        return response.Error(c, fiber.StatusNotFound, err)
    }
    return response.JSON(c, fiber.StatusOK, lead)
}

func (h *Handler) ScoreLead(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    leadID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    if err := h.services.CRM.ScoreLead(tenantUUID, leadID); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusOK, map[string]string{"status": "scored"})
}

func (h *Handler) CreateIntegration(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    var payload dto.IntegrationRequest
    if err := h.validate(c, &payload); err != nil {
        return err
    }
    configBytes, err := json.Marshal(payload.Config)
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    integration := &models.Integration{
        BaseModel: models.BaseModel{ID: uuid.New(), TenantID: tenantUUID},
        Name:      payload.Name,
        Provider:  payload.Provider,
        Config:     datatypes.JSON(configBytes),
        Enabled:   true,
    }
    if err := h.services.Integration.CreateIntegration(integration); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusCreated, integration)
}

func (h *Handler) ListIntegrations(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    integrations, err := h.services.Integration.ListIntegrations(tenantUUID)
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusOK, integrations)
}

func (h *Handler) GetPlans(c *fiber.Ctx) error {
    return response.JSON(c, fiber.StatusOK, h.services.Billing.GetPlans())
}

func (h *Handler) CreateSubscription(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    var payload dto.SubscriptionRequest
    if err := h.validate(c, &payload); err != nil {
        return err
    }
    subscription := &models.Subscription{
        BaseModel: models.BaseModel{ID: uuid.New(), TenantID: tenantUUID},
        Plan:      payload.Plan,
        Status:    "active",
        RenewAt:   time.Now().Add(30 * 24 * time.Hour),
        Metadata:  datatypes.JSON(json.RawMessage(`{}`)),
    }
    if err := h.services.Billing.CreateSubscription(subscription); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusCreated, subscription)
}

func (h *Handler) CreateInvoice(c *fiber.Ctx) error {
    tenantID := c.Locals("tenant_id").(string)
    tenantUUID, _ := uuid.Parse(tenantID)
    var payload dto.CreateInvoiceRequest
    if err := h.validate(c, &payload); err != nil {
        return err
    }
    invoice := &models.Invoice{
        BaseModel:   models.BaseModel{ID: uuid.New(), TenantID: tenantUUID},
        AmountCents: payload.AmountCents,
        Currency:    payload.Currency,
        Status:      "pending",
        DueAt:       time.Now().Add(7 * 24 * time.Hour),
    }
    if err := h.services.Billing.CreateInvoice(invoice); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusCreated, invoice)
}

func (h *Handler) ListTenants(c *fiber.Ctx) error {
    tenants, err := h.services.Admin.ListTenants()
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusOK, tenants)
}

func (h *Handler) SuspendTenant(c *fiber.Ctx) error {
    tenantID, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.Error(c, fiber.StatusBadRequest, err)
    }
    if err := h.services.Admin.SuspendTenant(tenantID); err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.Message(c, fiber.StatusOK, "tenant suspended")
}

func (h *Handler) ListAuditLogs(c *fiber.Ctx) error {
    logs, err := h.services.Admin.ListAuditLogs()
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusOK, logs)
}

func (h *Handler) GetAnalytics(c *fiber.Ctx) error {
    summary, err := h.services.Analytics.Summary()
    if err != nil {
        return response.Error(c, fiber.StatusInternalServerError, err)
    }
    return response.JSON(c, fiber.StatusOK, summary)
}

func (h *Handler) StartScheduler(client *asynq.Client) error {
    if h.scheduler != nil {
        return nil
    }
    scheduler := asynq.NewScheduler(asynq.RedisClientOpt{Addr: h.cfg.RedisURL}, nil)
    h.scheduler = scheduler
    return nil
}
