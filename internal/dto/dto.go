package dto

import "github.com/google/uuid"

type RegisterRequest struct {
    TenantName string `json:"tenant_name" validate:"required,min=3,max=64"`
    Name       string `json:"name" validate:"required,min=2,max=64"`
    Email      string `json:"email" validate:"required,email"`
    Password   string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
    TenantID string `json:"tenant_id" validate:"required,uuid"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type RefreshRequest struct {
    RefreshToken string `json:"refresh_token" validate:"required"`
}

type TenantRequest struct {
    Name string `json:"name" validate:"required,min=3"`
}

type AgentRequest struct {
    Name        string `json:"name" validate:"required,min=2"`
    Type        string `json:"type" validate:"required"`
    Description string `json:"description"`
    Prompt      string `json:"prompt"`
}

type AgentRunRequest struct {
    Input string `json:"input" validate:"required"`
}

type WorkflowRequest struct {
    Name        string `json:"name" validate:"required"`
    Description string `json:"description"`
    Trigger     string `json:"trigger"`
    Actions     any    `json:"actions" validate:"required"`
    Schedule    string `json:"schedule"`
}

type WorkflowRunRequest struct {
    Input string `json:"input" validate:"required"`
}

type LeadRequest struct {
    Name    string `json:"name" validate:"required"`
    Email   string `json:"email" validate:"required,email"`
    Company string `json:"company"`
}

type IntegrationRequest struct {
    Name     string      `json:"name" validate:"required"`
    Provider string      `json:"provider" validate:"required"`
    Config   interface{} `json:"config"`
}

type SubscriptionRequest struct {
    Plan string `json:"plan" validate:"required,oneof=free pro enterprise"`
}

type CreateInvoiceRequest struct {
    AmountCents int64  `json:"amount_cents" validate:"required,gt=0"`
    Currency    string `json:"currency" validate:"required"`
}

type InviteUserRequest struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required"`
    Role  string `json:"role" validate:"required,oneof=owner admin member"`
}

type IDRequest struct {
    ID uuid.UUID `json:"id" validate:"required,uuid"`
}
