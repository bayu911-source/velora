package services

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/google/uuid"
    "github.com/golang-jwt/jwt/v5"
    "github.com/hibiken/asynq"
    "go.uber.org/zap"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/datatypes"

    "velora/config"
    "velora/internal/models"
    "velora/internal/repositories"
    "velora/internal/workers"
)

const (
    DefaultPlan          = "free"
    RoleOwner            = "owner"
    RoleAdmin            = "admin"
    RoleMember           = "member"
    WorkflowTaskType     = workers.TypeWorkflowRun
    EmailTaskType        = workers.TypeEmailSend
    AIProcessingTaskType = workers.TypeAIProcess
)

// TokenClaims stores JWT claims for auth middleware.
type TokenClaims struct {
    UserID   string `json:"user_id"`
    TenantID string `json:"tenant_id"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

// AIService sends prompts to Gemini and returns generated text.
type AIService struct {
    cfg        *config.Config
    httpClient *http.Client
}

func NewAIService(cfg *config.Config) *AIService {
    return &AIService{
        cfg: cfg,
        httpClient: &http.Client{Timeout: 20 * time.Second},
    }
}

func (s *AIService) Generate(prompt string) (string, error) {
    if strings.TrimSpace(prompt) == "" {
        return "", errors.New("prompt cannot be empty")
    }

    requestBody := map[string]any{
        "prompt": map[string]any{
            "text": prompt,
        },
        "temperature": 0.7,
        "maxOutputTokens": 512,
    }
    bodyBytes, err := json.Marshal(requestBody)
    if err != nil {
        return "", err
    }

    requestURL := fmt.Sprintf("%s/v1beta2/models/text-bison-001:generate", strings.TrimRight(s.cfg.GeminiAPIURL, "/"))
    req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, requestURL, strings.NewReader(string(bodyBytes)))
    if err != nil {
        return "", err
    }
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.cfg.GeminiAPIKey))
    req.Header.Set("Content-Type", "application/json")

    res, err := s.httpClient.Do(req)
    if err != nil {
        return "", err
    }
    defer res.Body.Close()

    if res.StatusCode >= 300 {
        return "", fmt.Errorf("gemini returned status %d", res.StatusCode)
    }

    var response struct {
        Candidates []struct {
            Output string `json:"output"`
        } `json:"candidates"`
    }
    if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
        return "", err
    }
    if len(response.Candidates) == 0 {
        return "", errors.New("gemini returned no output")
    }
    return response.Candidates[0].Output, nil
}

// AuthService supports registration, login, and refresh token flows.
type AuthService struct {
    cfg  *config.Config
    repo *repositories.Repository
}

func NewAuthService(cfg *config.Config, repo *repositories.Repository) *AuthService {
    return &AuthService{cfg: cfg, repo: repo}
}

func (s *AuthService) RegisterTenant(tenantName, userName, email, password string) (*models.Tenant, *models.User, string, string, error) {
    tenantID := uuid.New()
    tenant := &models.Tenant{
        ID:      tenantID,
        Name:    tenantName,
        Slug:    strings.ToLower(strings.ReplaceAll(tenantName, " ", "-")),
        Plan:    DefaultPlan,
        Status:  "active",
        Branding: datatypes.JSON(json.RawMessage(`{"theme":"default"}`)),
    }
    if err := s.repo.CreateTenant(tenant); err != nil {
        return nil, nil, "", "", err
    }

    passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, nil, "", "", err
    }

    user := &models.User{
        BaseModel:    models.BaseModel{ID: uuid.New(), TenantID: tenant.ID},
        Email:        strings.ToLower(email),
        Name:         userName,
        Role:         RoleOwner,
        PasswordHash: string(passwordHash),
        IsActive:     true,
    }
    if err := s.repo.CreateUser(user); err != nil {
        return nil, nil, "", "", err
    }

    access, refresh, err := s.GenerateSessionTokens(user)
    if err != nil {
        return nil, nil, "", "", err
    }
    return tenant, user, access, refresh, nil
}

func (s *AuthService) Login(tenantID uuid.UUID, email, password string) (*models.User, string, string, error) {
    user, err := s.repo.FindUserByEmail(tenantID, strings.ToLower(email))
    if err != nil {
        return nil, "", "", err
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return nil, "", "", err
    }

    user.LastLogin = time.Now().UTC()
    if err := s.repo.UpdateUser(user); err != nil {
        return nil, "", "", err
    }

    access, refresh, err := s.GenerateSessionTokens(user)
    if err != nil {
        return nil, "", "", err
    }
    return user, access, refresh, nil
}

func (s *AuthService) GenerateSessionTokens(user *models.User) (string, string, error) {
    now := time.Now().UTC()
    accessClaims := TokenClaims{
        UserID:   user.ID.String(),
        TenantID: user.TenantID.String(),
        Role:     user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.AccessTokenTTL)),
            Subject:   user.ID.String(),
        },
    }

    refreshClaims := TokenClaims{
        UserID:   user.ID.String(),
        TenantID: user.TenantID.String(),
        Role:     user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.RefreshTokenTTL)),
            Subject:   user.ID.String(),
        },
    }

    accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.cfg.JWTSecret))
    if err != nil {
        return "", "", err
    }

    refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.cfg.RefreshSecret))
    if err != nil {
        return "", "", err
    }
    return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (string, string, error) {
    token, err := jwt.ParseWithClaims(refreshToken, &TokenClaims{}, func(token *jwt.Token) (any, error) {
        if token.Method != jwt.SigningMethodHS256 {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(s.cfg.RefreshSecret), nil
    })
    if err != nil {
        return "", "", err
    }

    claims, ok := token.Claims.(*TokenClaims)
    if !ok || !token.Valid {
        return "", "", errors.New("invalid refresh token")
    }

    userID, err := uuid.Parse(claims.UserID)
    if err != nil {
        return "", "", err
    }
    tenantID, err := uuid.Parse(claims.TenantID)
    if err != nil {
        return "", "", err
    }

    user, err := s.repo.FindUserByID(tenantID, userID)
    if err != nil {
        return "", "", err
    }

    return s.GenerateSessionTokens(user)
}

func (s *AuthService) InviteUser(tenantID uuid.UUID, email, name, role string) (*models.User, error) {
    user := &models.User{
        BaseModel: models.BaseModel{ID: uuid.New(), TenantID: tenantID},
        Email:     strings.ToLower(email),
        Name:      name,
        Role:      role,
        IsActive:  true,
    }
    if err := s.repo.CreateUser(user); err != nil {
        return nil, err
    }
    return user, nil
}

// TenantService manages tenant workspace operations.
type TenantService struct {
    repo *repositories.Repository
}

func NewTenantService(repo *repositories.Repository) *TenantService {
    return &TenantService{repo: repo}
}

func (s *TenantService) CreateTenant(name string) (*models.Tenant, error) {
    tenant := &models.Tenant{
        ID:       uuid.New(),
        Name:     name,
        Slug:     strings.ToLower(strings.ReplaceAll(name, " ", "-")),
        Plan:     DefaultPlan,
        Status:   "active",
        Branding: datatypes.JSON(json.RawMessage(`{"theme":"default"}`)),
    }
    if err := s.repo.CreateTenant(tenant); err != nil {
        return nil, err
    }
    return tenant, nil
}

func (s *TenantService) GetTenant(id uuid.UUID) (*models.Tenant, error) {
    return s.repo.GetTenant(id)
}

func (s *TenantService) SuspendTenant(id uuid.UUID) error {
    return s.repo.SuspendTenant(id)
}

// AgentService orchestrates agent CRUD and prompt generation.
type AgentService struct {
    repo *repositories.Repository
    ai   *AIService
}

func NewAgentService(repo *repositories.Repository, ai *AIService) *AgentService {
    return &AgentService{repo: repo, ai: ai}
}

func (s *AgentService) CreateAgent(agent *models.Agent) error {
    return s.repo.CreateAgent(agent)
}

func (s *AgentService) ListAgents(tenantID uuid.UUID) ([]models.Agent, error) {
    return s.repo.ListAgents(tenantID)
}

func (s *AgentService) GetAgent(tenantID uuid.UUID, id uuid.UUID) (*models.Agent, error) {
    return s.repo.GetAgent(tenantID, id)
}

func (s *AgentService) UpdateAgent(agent *models.Agent) error {
    return s.repo.UpdateAgent(agent)
}

func (s *AgentService) DeleteAgent(tenantID uuid.UUID, id uuid.UUID) error {
    return s.repo.DeleteAgent(tenantID, id)
}

func (s *AgentService) RunAgent(tenantID uuid.UUID, agentID uuid.UUID, input string) (string, error) {
    agent, err := s.repo.GetAgent(tenantID, agentID)
    if err != nil {
        return "", err
    }
    prompt := strings.TrimSpace(agent.Prompt)
    if prompt == "" {
        prompt = fmt.Sprintf("Run agent %s: %s", agent.Name, input)
    }
    fullPrompt := fmt.Sprintf("%s\n\nInput:\n%s", prompt, input)
    return s.ai.Generate(fullPrompt)
}

// WorkflowService manages automation definitions and execution.
type WorkflowService struct {
    repo       *repositories.Repository
    agent      *AgentService
    client     *asynq.Client
    logger     *zap.Logger
}

func NewWorkflowService(repo *repositories.Repository, agent *AgentService, client *asynq.Client, logger *zap.Logger) *WorkflowService {
    return &WorkflowService{repo: repo, agent: agent, client: client, logger: logger}
}

func (s *WorkflowService) CreateWorkflow(workflow *models.Workflow) error {
    return s.repo.CreateWorkflow(workflow)
}

func (s *WorkflowService) ListWorkflows(tenantID uuid.UUID) ([]models.Workflow, error) {
    return s.repo.ListWorkflows(tenantID)
}

func (s *WorkflowService) GetWorkflow(tenantID uuid.UUID, id uuid.UUID) (*models.Workflow, error) {
    return s.repo.GetWorkflow(tenantID, id)
}

func (s *WorkflowService) DeleteWorkflow(tenantID uuid.UUID, id uuid.UUID) error {
    return s.repo.DeleteWorkflow(tenantID, id)
}

func (s *WorkflowService) ScheduleWorkflow(workflow *models.Workflow) error {
    if workflow.Schedule == "" {
        return nil
    }
    // schedule is persisted and may be picked up by the worker scheduler
    return s.repo.CreateWorkflow(workflow)
}

func (s *WorkflowService) EnqueueWorkflowRun(tenantID uuid.UUID, workflowID uuid.UUID, input string) error {
    payload := map[string]any{"workflow_id": workflowID.String(), "tenant_id": tenantID.String(), "input": input}
    task, err := workers.NewTask(workers.TypeWorkflowRun, payload)
    if err != nil {
        return err
    }
    _, err = s.client.Enqueue(task, asynq.MaxRetry(3), asynq.ProcessIn(1*time.Second))
    return err
}

func (s *WorkflowService) ExecuteWorkflow(ctx context.Context, tenantID uuid.UUID, workflow *models.Workflow, input string) (string, error) {
    if len(workflow.Actions) == 0 {
        return "", errors.New("workflow contains no actions")
    }

    var actions []struct {
        Name    string `json:"name"`
        AgentID string `json:"agent_id"`
        Prompt  string `json:"prompt"`
    }
    if err := json.Unmarshal(workflow.Actions, &actions); err != nil {
        return "", err
    }

    var output string
    for _, action := range actions {
        if action.AgentID == "" {
            return "", errors.New("workflow action is missing agent_id")
        }
        agentID, err := uuid.Parse(action.AgentID)
        if err != nil {
            return "", err
        }
        text := action.Prompt
        if text == "" {
            text = output
        }
        output, err = s.agent.RunAgent(tenantID, agentID, text)
        if err != nil {
            return "", err
        }
    }

    run := &models.WorkflowRun{
        BaseModel:  models.BaseModel{ID: uuid.New(), TenantID: tenantID},
        WorkflowID: workflow.ID,
        Status:     "completed",
        Input:      datatypes.JSON(json.RawMessage(fmt.Sprintf(`{"input":"%s"}`, input))),
        Output:     datatypes.JSON(json.RawMessage(fmt.Sprintf(`{"output":"%s"}`, output))),
    }
    _ = s.repo.CreateWorkflowRun(run)
    return output, nil
}

func (s *WorkflowService) HandleWorkflowTask(ctx context.Context, task *asynq.Task) error {
    var payload struct {
        WorkflowID string `json:"workflow_id"`
        TenantID   string `json:"tenant_id"`
        Input      string `json:"input"`
    }
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return err
    }
    tenantID, err := uuid.Parse(payload.TenantID)
    if err != nil {
        return err
    }
    workflowID, err := uuid.Parse(payload.WorkflowID)
    if err != nil {
        return err
    }
    workflow, err := s.repo.GetWorkflow(tenantID, workflowID)
    if err != nil {
        return err
    }
    if workflow == nil {
        return errors.New("workflow not found")
    }
    _, err = s.ExecuteWorkflow(ctx, tenantID, workflow, payload.Input)
    return err
}

// CRMService exposes lead management and scoring.
type CRMService struct {
    repo *repositories.Repository
    ai   *AIService
}

func NewCRMService(repo *repositories.Repository, ai *AIService) *CRMService {
    return &CRMService{repo: repo, ai: ai}
}

func (s *CRMService) CreateLead(lead *models.Lead) error {
    return s.repo.CreateLead(lead)
}

func (s *CRMService) ListLeads(tenantID uuid.UUID) ([]models.Lead, error) {
    return s.repo.ListLeads(tenantID)
}

func (s *CRMService) GetLead(tenantID uuid.UUID, id uuid.UUID) (*models.Lead, error) {
    return s.repo.GetLead(tenantID, id)
}

func (s *CRMService) ScoreLead(tenantID uuid.UUID, leadID uuid.UUID) error {
    lead, err := s.repo.GetLead(tenantID, leadID)
    if err != nil {
        return err
    }
    prompt := fmt.Sprintf("Provide a lead score from 1 to 100 for the following lead: Name=%s, Company=%s, Email=%s", lead.Name, lead.Company, lead.Email)
    text, err := s.ai.Generate(prompt)
    if err != nil {
        return err
    }
    score := 50
    if parsed, pErr := fmt.Sscanf(text, "%d", &score); pErr == nil && parsed == 1 {
        lead.Score = score
    } else {
        lead.Score = 50
    }
    return s.repo.UpdateLead(lead)
}

// IntegrationService stores third-party service connectors.
type IntegrationService struct {
    repo *repositories.Repository
}

func NewIntegrationService(repo *repositories.Repository) *IntegrationService {
    return &IntegrationService{repo: repo}
}

func (s *IntegrationService) CreateIntegration(integration *models.Integration) error {
    return s.repo.CreateIntegration(integration)
}

func (s *IntegrationService) ListIntegrations(tenantID uuid.UUID) ([]models.Integration, error) {
    return s.repo.ListIntegrations(tenantID)
}

// BillingService manages plans, invoices, and subscriptions.
type BillingService struct {
    repo *repositories.Repository
}

func NewBillingService(repo *repositories.Repository) *BillingService {
    return &BillingService{repo: repo}
}

func (s *BillingService) GetPlans() []map[string]any {
    return []map[string]any{
        {"id": "free", "name": "Free", "limit": 100, "price": 0},
        {"id": "pro", "name": "Pro", "limit": 2000, "price": 49},
        {"id": "enterprise", "name": "Enterprise", "limit": 10000, "price": 249},
    }
}

func (s *BillingService) CreateSubscription(subscription *models.Subscription) error {
    return s.repo.CreateSubscription(subscription)
}

func (s *BillingService) CreateInvoice(invoice *models.Invoice) error {
    return s.repo.CreateInvoice(invoice)
}

func (s *BillingService) EnforceUsageLimit(tenantID uuid.UUID) error {
    subscription, err := s.repo.GetSubscription(tenantID)
    if err != nil {
        return err
    }
    if subscription.Plan == "free" {
        count, err := s.repo.CountUsage(tenantID, "workflow_execution")
        if err != nil {
            return err
        }
        if count > 500 {
            return errors.New("free plan usage limit reached")
        }
    }
    return nil
}

// AdminService provides global tenant management.
type AdminService struct {
    repo *repositories.Repository
}

func NewAdminService(repo *repositories.Repository) *AdminService {
    return &AdminService{repo: repo}
}

func (s *AdminService) ListTenants() ([]models.Tenant, error) {
    return s.repo.ListTenants()
}

func (s *AdminService) SuspendTenant(tenantID uuid.UUID) error {
    return s.repo.SuspendTenant(tenantID)
}

func (s *AdminService) ListAuditLogs() ([]models.AuditLog, error) {
    return s.repo.ListAuditLogs()
}

// AnalyticsService tracks platform usage for dashboards.
type AnalyticsService struct {
    repo *repositories.Repository
}

func NewAnalyticsService(repo *repositories.Repository) *AnalyticsService {
    return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) Summary() (map[string]any, error) {
    tenants, err := s.repo.CountRecords(&models.Tenant{})
    if err != nil {
        return nil, err
    }
    workflows, err := s.repo.CountRecords(&models.Workflow{})
    if err != nil {
        return nil, err
    }
    leads, err := s.repo.CountRecords(&models.Lead{})
    if err != nil {
        return nil, err
    }
    return map[string]any{
        "tenants":    tenants,
        "workflows":  workflows,
        "leads":      leads,
        "generated_at": time.Now().UTC(),
    }, nil
}

// Services aggregates all domain services used by handlers.
type Services struct {
    Auth        *AuthService
    Tenant      *TenantService
    Agent       *AgentService
    Workflow    *WorkflowService
    CRM         *CRMService
    Integration *IntegrationService
    Billing     *BillingService
    Admin       *AdminService
    Analytics   *AnalyticsService
}

func NewServices(cfg *config.Config, repo *repositories.Repository, logger *zap.Logger) *Services {
    aiService := NewAIService(cfg)
    asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisURL})
    agentService := NewAgentService(repo, aiService)
    workflowService := NewWorkflowService(repo, agentService, asynqClient, logger)
    return &Services{
        Auth:        NewAuthService(cfg, repo),
        Tenant:      NewTenantService(repo),
        Agent:       agentService,
        Workflow:    workflowService,
        CRM:         NewCRMService(repo, aiService),
        Integration: NewIntegrationService(repo),
        Billing:     NewBillingService(repo),
        Admin:       NewAdminService(repo),
        Analytics:   NewAnalyticsService(repo),
    }
}
