package workers

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hibiken/asynq"
    "github.com/google/uuid"
    "go.uber.org/zap"
    "velora/internal/models"
)

// WorkflowExecutor defines the operations the worker needs from the workflow service.
type WorkflowExecutor interface {
    GetWorkflow(tenantID uuid.UUID, id uuid.UUID) (*models.Workflow, error)
    ExecuteWorkflow(ctx context.Context, tenantID uuid.UUID, workflow *models.Workflow, input string) (string, error)
}

// WorkerServer wraps Asynq server and task processing.
type WorkerServer struct {
    server  *asynq.Server
    handler WorkflowExecutor
    logger  *zap.Logger
}

func NewWorkerServer(cfgAddr string, workflowService WorkflowExecutor, logger *zap.Logger) *WorkerServer {
    server := asynq.NewServer(asynq.RedisClientOpt{Addr: cfgAddr}, asynq.Config{Concurrency: 10})
    return &WorkerServer{server: server, handler: workflowService, logger: logger}
}

func (w *WorkerServer) Start(ctx context.Context) error {
    mux := asynq.NewServeMux()
    mux.HandleFunc(TypeWorkflowRun, w.handleWorkflowRun)
    mux.HandleFunc(TypeEmailSend, w.handleEmailSend)
    mux.HandleFunc(TypeAIProcess, w.handleAIProcess)

    w.logger.Info("starting worker server")
    return w.server.Run(mux)
}

func (w *WorkerServer) handleWorkflowRun(ctx context.Context, task *asynq.Task) error {
    var payload TaskPayload
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
    workflow, err := w.handler.GetWorkflow(tenantID, workflowID)
    if err != nil {
        return err
    }
    if workflow == nil {
        return fmt.Errorf("workflow %s not found", workflowID)
    }
    _, err = w.handler.ExecuteWorkflow(ctx, tenantID, workflow, payload.Input)
    return err
}

func (w *WorkerServer) handleEmailSend(ctx context.Context, task *asynq.Task) error {
    var payload TaskPayload
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return err
    }
    w.logger.Info("mock email sent", zap.String("subject", payload.Subject), zap.String("body", payload.Body))
    return nil
}

func (w *WorkerServer) handleAIProcess(ctx context.Context, task *asynq.Task) error {
    var payload TaskPayload
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return err
    }
    w.logger.Info("processing ai task", zap.String("input", payload.Input))
    return nil
}
