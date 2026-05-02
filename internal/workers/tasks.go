package workers

import (
    "encoding/json"

    "github.com/google/uuid"
    "github.com/hibiken/asynq"
)

const (
    TypeWorkflowRun = "workflow:run"
    TypeEmailSend   = "email:send"
    TypeAIProcess   = "ai:process"
)

type TaskPayload struct {
    WorkflowID string `json:"workflow_id,omitempty"`
    TenantID   string `json:"tenant_id,omitempty"`
    Input      string `json:"input,omitempty"`
    Subject    string `json:"subject,omitempty"`
    Body       string `json:"body,omitempty"`
}

func NewTask(taskType string, payload any) (*asynq.Task, error) {
    body, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }
    return asynq.NewTask(taskType, body), nil
}

func ValidateUUID(id string) error {
    if _, err := uuid.Parse(id); err != nil {
        return err
    }
    return nil
}
