package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/datatypes"
    "gorm.io/gorm"
)

// BaseModel defines fields shared by all tenant-scoped records.
type BaseModel struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    TenantID  uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
    if m.ID == uuid.Nil {
        m.ID = uuid.New()
    }
    return nil
}

// Tenant represents an isolated customer workspace.
type Tenant struct {
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
    Name      string         `gorm:"not null" json:"name"`
    Slug      string         `gorm:"uniqueIndex;not null" json:"slug"`
    Plan      string         `gorm:"not null;default:'free'" json:"plan"`
    Branding  datatypes.JSON `gorm:"type:jsonb" json:"branding"`
    Status    string         `gorm:"not null;default:'active'" json:"status"`
    OwnerID   uuid.UUID      `gorm:"type:uuid;index" json:"owner_id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
}

func (t *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
    if t.ID == uuid.Nil {
        t.ID = uuid.New()
    }
    return nil
}

// User is an authenticated identity scoped to a tenant.
type User struct {
    BaseModel
    Email        string    `gorm:"not null;index:idx_user_email,unique,priority:1" json:"email"`
    Name         string    `gorm:"not null" json:"name"`
    Role         string    `gorm:"not null;default:'member'" json:"role"`
    PasswordHash string    `gorm:"not null" json:"-"`
    IsActive     bool      `gorm:"default:true" json:"is_active"`
    LastLogin    time.Time `json:"last_login"`
    InviteToken  string    `gorm:"index" json:"invite_token,omitempty"`
}

// APIKey allows scoped token-based access for automation.
type APIKey struct {
    BaseModel
    Name       string    `gorm:"not null" json:"name"`
    Key        string    `gorm:"not null;uniqueIndex" json:"key"`
    Scopes     string    `gorm:"type:text;default:'read,write'" json:"scopes"`
    Active     bool      `gorm:"default:true" json:"active"`
    LastUsedAt time.Time `json:"last_used_at"`
}

// Agent represents a reusable AI automation agent.
type Agent struct {
    BaseModel
    Name        string         `gorm:"not null" json:"name"`
    Type        string         `gorm:"not null" json:"type"`
    Description string         `json:"description"`
    Prompt      string         `gorm:"type:text" json:"prompt"`
    Config      datatypes.JSON `gorm:"type:jsonb" json:"config"`
    Active      bool           `gorm:"default:true" json:"active"`
}

// Workflow stores automation flow definitions with trigger and actions.
type Workflow struct {
    BaseModel
    Name        string         `gorm:"not null" json:"name"`
    Description string         `json:"type:text" json:"description"`
    Trigger     string         `gorm:"type:text" json:"trigger"`
    Actions     datatypes.JSON `gorm:"type:jsonb" json:"actions"`
    Schedule    string         `json:"schedule"`
    Status      string         `gorm:"not null;default:'draft'" json:"status"`
}

// WorkflowRun captures asynchronous workflow execution history.
type WorkflowRun struct {
    BaseModel
    WorkflowID uuid.UUID      `gorm:"type:uuid;index;not null" json:"workflow_id"`
    Status     string         `gorm:"not null;default:'pending'" json:"status"`
    Input      datatypes.JSON `gorm:"type:jsonb" json:"input"`
    Output     datatypes.JSON `gorm:"type:jsonb" json:"output"`
    Error      string         `gorm:"type:text" json:"error"`
}

// Lead is a record for potential customers and opportunity tracking.
type Lead struct {
    BaseModel
    Name        string    `gorm:"not null" json:"name"`
    Email       string    `gorm:"index" json:"email"`
    Company     string    `json:"company"`
    Status      string    `gorm:"not null;default:'new'" json:"status"`
    Score       int       `gorm:"default:0" json:"score"`
    OwnerID     uuid.UUID `gorm:"type:uuid;index" json:"owner_id"`
    Pipeline    string    `gorm:"default:'qualification'" json:"pipeline"`
}

// Contact stores people linked to leads and accounts.
type Contact struct {
    BaseModel
    LeadID uuid.UUID `gorm:"type:uuid;index;not null" json:"lead_id"`
    Name   string    `gorm:"not null" json:"name"`
    Email  string    `gorm:"json:"email"`
    Phone  string    `json:"phone"`
    Notes  string    `gorm:"type:text" json:"notes"`
}

// Note stores notes attached to leads, contacts, or workflows.
type Note struct {
    BaseModel
    LeadID uuid.UUID `gorm:"type:uuid;index" json:"lead_id"`
    Body   string    `gorm:"type:text;not null" json:"body"`
    Author string    `json:"author"`
}

// Integration holds a connection to a third-party or mock service.
type Integration struct {
    BaseModel
    Name    string         `gorm:"not null" json:"name"`
    Provider string        `gorm:"not null" json:"provider"`
    Config  datatypes.JSON `gorm:"type:jsonb" json:"config"`
    Enabled bool           `gorm:"default:true" json:"enabled"`
}

// Subscription tracks tenant plan and billing status.
type Subscription struct {
    BaseModel
    TenantID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
    Plan        string         `gorm:"not null" json:"plan"`
    Status      string         `gorm:"not null;default:'active'" json:"status"`
    RenewAt     time.Time      `json:"renew_at"`
    Metadata    datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
}

// Invoice records charges and payment state.
type Invoice struct {
    BaseModel
    TenantID   uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
    AmountCents int64    `gorm:"not null" json:"amount_cents"`
    Currency    string   `gorm:"not null;default:'USD'" json:"currency"`
    Status      string   `gorm:"not null;default:'pending'" json:"status"`
    DueAt       time.Time `json:"due_at"`
}

// UsageRecord tracks API and workflow usage per tenant.
type UsageRecord struct {
    BaseModel
    TenantID uuid.UUID      `gorm:"type:uuid;index;not null" json:"tenant_id"`
    Category string         `gorm:"not null" json:"category"`
    Quantity int64          `gorm:"not null" json:"quantity"`
    Meta     datatypes.JSON `gorm:"type:jsonb" json:"meta"`
}

// Partner represents affiliate or channel partners.
type Partner struct {
    BaseModel
    Company string `gorm:"not null" json:"company"`
    Code    string `gorm:"not null;uniqueIndex" json:"code"`
    Active  bool   `gorm:"default:true" json:"active"`
}

// Referral records affiliate referrals.
type Referral struct {
    BaseModel
    PartnerID uuid.UUID `gorm:"type:uuid;index;not null" json:"partner_id"`
    TenantID  uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
    Source    string    `json:"source"`
    Reward    int64     `json:"reward"`
}

// AuditLog stores security and operations audit trails.
type AuditLog struct {
    BaseModel
    ActorID    uuid.UUID `gorm:"type:uuid;index" json:"actor_id"`
    Entity     string    `gorm:"not null" json:"entity"`
    Action     string    `gorm:"not null" json:"action"`
    ResourceID uuid.UUID `gorm:"type:uuid;index" json:"resource_id"`
    Message    string    `gorm:"type:text" json:"message"`
}

// Notification tracks user-facing notifications.
type Notification struct {
    BaseModel
    TenantID uuid.UUID `gorm:"type:uuid;index;not null" json:"tenant_id"`
    Title    string    `gorm:"not null" json:"title"`
    Body     string    `gorm:"type:text;not null" json:"body"`
    Read     bool      `gorm:"default:false" json:"read"`
}
