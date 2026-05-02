package repositories

import (
    "errors"

    "github.com/google/uuid"
    "gorm.io/gorm"

    "velora/internal/models"
)

// Repository provides tenant-aware access to the data layer.
type Repository struct {
    db *gorm.DB
}

// NewRepository initializes a new repository instance.
func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) tenantQuery(tenantID uuid.UUID) *gorm.DB {
    return r.db.Where("tenant_id = ?", tenantID)
}

func (r *Repository) CreateTenant(tenant *models.Tenant) error {
    return r.db.Create(tenant).Error
}

func (r *Repository) GetTenant(id uuid.UUID) (*models.Tenant, error) {
    var tenant models.Tenant
    if err := r.db.First(&tenant, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &tenant, nil
}

func (r *Repository) ListTenants() ([]models.Tenant, error) {
    var tenants []models.Tenant
    if err := r.db.Order("created_at desc").Find(&tenants).Error; err != nil {
        return nil, err
    }
    return tenants, nil
}

func (r *Repository) SuspendTenant(id uuid.UUID) error {
    return r.db.Model(&models.Tenant{}).Where("id = ?", id).Updates(map[string]interface{}{"status": "suspended"}).Error
}

func (r *Repository) CreateUser(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *Repository) UpdateUser(user *models.User) error {
    return r.db.Save(user).Error
}

func (r *Repository) FindUserByEmail(tenantID uuid.UUID, email string) (*models.User, error) {
    var user models.User
    if err := r.tenantQuery(tenantID).Where("email = ?", email).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *Repository) FindUserByID(tenantID uuid.UUID, userID uuid.UUID) (*models.User, error) {
    var user models.User
    if err := r.tenantQuery(tenantID).Where("id = ?", userID).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *Repository) CreateAPIKey(apiKey *models.APIKey) error {
    return r.db.Create(apiKey).Error
}

func (r *Repository) FindAPIKey(key string) (*models.APIKey, error) {
    var apiKey models.APIKey
    if err := r.db.Where("key = ? AND active = true", key).First(&apiKey).Error; err != nil {
        return nil, err
    }
    return &apiKey, nil
}

func (r *Repository) RevokeAPIKey(tenantID uuid.UUID, id uuid.UUID) error {
    return r.tenantQuery(tenantID).Model(&models.APIKey{}).Where("id = ?", id).Update("active", false).Error
}

func (r *Repository) CreateAgent(agent *models.Agent) error {
    return r.db.Create(agent).Error
}

func (r *Repository) ListAgents(tenantID uuid.UUID) ([]models.Agent, error) {
    var agents []models.Agent
    if err := r.tenantQuery(tenantID).Order("created_at desc").Find(&agents).Error; err != nil {
        return nil, err
    }
    return agents, nil
}

func (r *Repository) GetAgent(tenantID uuid.UUID, id uuid.UUID) (*models.Agent, error) {
    var agent models.Agent
    if err := r.tenantQuery(tenantID).Where("id = ?", id).First(&agent).Error; err != nil {
        return nil, err
    }
    return &agent, nil
}

func (r *Repository) UpdateAgent(agent *models.Agent) error {
    return r.db.Save(agent).Error
}

func (r *Repository) DeleteAgent(tenantID uuid.UUID, id uuid.UUID) error {
    return r.tenantQuery(tenantID).Where("id = ?", id).Delete(&models.Agent{}).Error
}

func (r *Repository) CreateWorkflow(workflow *models.Workflow) error {
    return r.db.Create(workflow).Error
}

func (r *Repository) ListWorkflows(tenantID uuid.UUID) ([]models.Workflow, error) {
    var workflows []models.Workflow
    if err := r.tenantQuery(tenantID).Order("created_at desc").Find(&workflows).Error; err != nil {
        return nil, err
    }
    return workflows, nil
}

func (r *Repository) GetWorkflow(tenantID uuid.UUID, id uuid.UUID) (*models.Workflow, error) {
    var workflow models.Workflow
    if err := r.tenantQuery(tenantID).Where("id = ?", id).First(&workflow).Error; err != nil {
        return nil, err
    }
    return &workflow, nil
}

func (r *Repository) DeleteWorkflow(tenantID uuid.UUID, id uuid.UUID) error {
    return r.tenantQuery(tenantID).Where("id = ?", id).Delete(&models.Workflow{}).Error
}

func (r *Repository) CreateWorkflowRun(run *models.WorkflowRun) error {
    return r.db.Create(run).Error
}

func (r *Repository) CreateLead(lead *models.Lead) error {
    return r.db.Create(lead).Error
}

func (r *Repository) ListLeads(tenantID uuid.UUID) ([]models.Lead, error) {
    var leads []models.Lead
    if err := r.tenantQuery(tenantID).Order("created_at desc").Find(&leads).Error; err != nil {
        return nil, err
    }
    return leads, nil
}

func (r *Repository) GetLead(tenantID uuid.UUID, id uuid.UUID) (*models.Lead, error) {
    var lead models.Lead
    if err := r.tenantQuery(tenantID).Where("id = ?", id).First(&lead).Error; err != nil {
        return nil, err
    }
    return &lead, nil
}

func (r *Repository) UpdateLead(lead *models.Lead) error {
    return r.db.Save(lead).Error
}

func (r *Repository) CreateNote(note *models.Note) error {
    return r.db.Create(note).Error
}

func (r *Repository) CreateIntegration(integration *models.Integration) error {
    return r.db.Create(integration).Error
}

func (r *Repository) ListIntegrations(tenantID uuid.UUID) ([]models.Integration, error) {
    var integrations []models.Integration
    if err := r.tenantQuery(tenantID).Order("created_at desc").Find(&integrations).Error; err != nil {
        return nil, err
    }
    return integrations, nil
}

func (r *Repository) CreateSubscription(subscription *models.Subscription) error {
    return r.db.Create(subscription).Error
}

func (r *Repository) CreateInvoice(invoice *models.Invoice) error {
    return r.db.Create(invoice).Error
}

func (r *Repository) CreateUsageRecord(record *models.UsageRecord) error {
    return r.db.Create(record).Error
}

func (r *Repository) CountUsage(tenantID uuid.UUID, category string) (int64, error) {
    var count int64
    if err := r.tenantQuery(tenantID).Model(&models.UsageRecord{}).Where("category = ?", category).Count(&count).Error; err != nil {
        return 0, err
    }
    return count, nil
}

func (r *Repository) CreateAuditLog(log *models.AuditLog) error {
    return r.db.Create(log).Error
}

func (r *Repository) ListAuditLogs() ([]models.AuditLog, error) {
    var logs []models.AuditLog
    if err := r.db.Order("created_at desc").Find(&logs).Error; err != nil {
        return nil, err
    }
    return logs, nil
}

func (r *Repository) CreateNotification(notification *models.Notification) error {
    return r.db.Create(notification).Error
}

func (r *Repository) CountRecords(model any) (int64, error) {
    var count int64
    if err := r.db.Model(model).Count(&count).Error; err != nil {
        return 0, err
    }
    return count, nil
}

func (r *Repository) EnsureTenantExists(tenantID uuid.UUID) error {
    var count int64
    if err := r.db.Model(&models.Tenant{}).Where("id = ?", tenantID).Count(&count).Error; err != nil {
        return err
    }
    if count == 0 {
        return errors.New("tenant not found")
    }
    return nil
}

func (r *Repository) GetSubscription(tenantID uuid.UUID) (*models.Subscription, error) {
    var sub models.Subscription
    if err := r.tenantQuery(tenantID).First(&sub).Error; err != nil {
        return nil, err
    }
    return &sub, nil
}
