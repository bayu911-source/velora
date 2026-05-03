package cmd

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/google/uuid"
    "github.com/spf13/cobra"
    "go.uber.org/zap"
    "gorm.io/datatypes"

    "velora/config"
    "velora/internal/database"
    "velora/internal/models"
    "velora/internal/repositories"
    "velora/internal/services"
)

type workflowAction struct {
    Name    string `json:"name" yaml:"name"`
    AgentID string `json:"agent_id" yaml:"agent_id"`
    Prompt  string `json:"prompt" yaml:"prompt"`
}

func NewWorkflowCmd() *cobra.Command {
    var tenantID string

    cmd := &cobra.Command{
        Use:   "workflow",
        Short: "Manage Velora workflows",
    }
    cmd.PersistentFlags().StringVarP(&tenantID, "tenant", "t", "", "Tenant UUID for workflow commands (or set TENANT_ID)")

    cmd.AddCommand(newWorkflowCreateCmd(&tenantID))
    cmd.AddCommand(newWorkflowListCmd(&tenantID))
    cmd.AddCommand(newWorkflowShowCmd(&tenantID))
    cmd.AddCommand(newWorkflowDeleteCmd(&tenantID))
    cmd.AddCommand(newWorkflowRunCmd(&tenantID))

    return cmd
}

func newWorkflowCreateCmd(tenantID *string) *cobra.Command {
    var description string
    var actionsFile string
    var agentIdentifiers string
    var actionPrompt string
    var status string

    cmd := &cobra.Command{
        Use:   "create <name> [description]",
        Short: "Create a new workflow",
        Args:  cobra.RangeArgs(1, 2),
        Run: func(cmd *cobra.Command, args []string) {
            if err := runWorkflowCreate(cmd, tenantID, args, description, actionsFile, agentIdentifiers, actionPrompt, status); err != nil {
                log.Fatalf("workflow create failed: %v", err)
            }
        },
    }
    cmd.Flags().StringVar(&description, "description", "", "Description for the workflow")
    cmd.Flags().StringVar(&actionsFile, "actions-file", "", "JSON file containing workflow actions")
    cmd.Flags().StringVar(&agentIdentifiers, "agents", "", "Comma-separated agent IDs or names for the workflow actions")
    cmd.Flags().StringVar(&actionPrompt, "prompt", "", "Prompt to use for each workflow action")
    cmd.Flags().StringVar(&status, "status", "active", "Workflow status (draft or active)")
    return cmd
}

func newWorkflowListCmd(tenantID *string) *cobra.Command {
    return &cobra.Command{
        Use:   "list",
        Short: "List workflows",
        Args:  cobra.NoArgs,
        Run: func(cmd *cobra.Command, args []string) {
            if err := runWorkflowList(cmd, tenantID); err != nil {
                log.Fatalf("workflow list failed: %v", err)
            }
        },
    }
}

func newWorkflowShowCmd(tenantID *string) *cobra.Command {
    return &cobra.Command{
        Use:   "show <workflow-id>",
        Short: "Show workflow details",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            if err := runWorkflowShow(cmd, tenantID, args[0]); err != nil {
                log.Fatalf("workflow show failed: %v", err)
            }
        },
    }
}

func newWorkflowDeleteCmd(tenantID *string) *cobra.Command {
    return &cobra.Command{
        Use:   "delete <workflow-id>",
        Short: "Delete a workflow",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            if err := runWorkflowDelete(cmd, tenantID, args[0]); err != nil {
                log.Fatalf("workflow delete failed: %v", err)
            }
        },
    }
}

func newWorkflowRunCmd(tenantID *string) *cobra.Command {
    var input string
    var async bool

    cmd := &cobra.Command{
        Use:   "run <workflow-id>",
        Short: "Execute a workflow",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            if err := runWorkflowRun(cmd, tenantID, args[0], input, async); err != nil {
                log.Fatalf("workflow run failed: %v", err)
            }
        },
    }
    cmd.Flags().StringVar(&input, "input", "", "Input payload for the workflow")
    cmd.Flags().BoolVar(&async, "async", false, "Enqueue workflow execution asynchronously")
    return cmd
}

func runWorkflowCreate(cmd *cobra.Command, tenantID *string, args []string, description, actionsFile, agentIdentifiers, actionPrompt, status string) error {
    tenantUUID, err := requireTenantID(*tenantID)
    if err != nil {
        return err
    }

    svc, err := loadServices()
    if err != nil {
        return err
    }

    name := args[0]
    if len(args) > 1 {
        description = args[1]
    }

    var actionsJSON []byte
    if actionsFile != "" {
        actionsJSON, err = os.ReadFile(actionsFile)
        if err != nil {
            return fmt.Errorf("failed to read actions file: %w", err)
        }
        if !json.Valid(actionsJSON) {
            return errors.New("actions file must contain valid JSON")
        }
    } else {
        if agentIdentifiers == "" {
            return errors.New("workflow create requires either --actions-file or --agents")
        }
        actions, err := buildActionsFromAgents(tenantUUID, svc.Agent, agentIdentifiers, actionPrompt)
        if err != nil {
            return err
        }
        actionsJSON, err = json.Marshal(actions)
        if err != nil {
            return fmt.Errorf("failed to encode workflow actions: %w", err)
        }
    }

    workflow := &models.Workflow{
        BaseModel:   models.BaseModel{ID: uuid.New(), TenantID: tenantUUID},
        Name:        name,
        Description: description,
        Trigger:     "manual",
        Actions:     datatypes.JSON(actionsJSON),
        Status:      status,
    }

    if err := svc.Workflow.CreateWorkflow(workflow); err != nil {
        return fmt.Errorf("failed to create workflow: %w", err)
    }
    fmt.Printf("Workflow created: %s\n", workflow.ID.String())
    return nil
}

func runWorkflowList(cmd *cobra.Command, tenantID *string) error {
    tenantUUID, err := requireTenantID(*tenantID)
    if err != nil {
        return err
    }

    svc, err := loadServices()
    if err != nil {
        return err
    }

    workflows, err := svc.Workflow.ListWorkflows(tenantUUID)
    if err != nil {
        return fmt.Errorf("failed to list workflows: %w", err)
    }

    if len(workflows) == 0 {
        fmt.Println("No workflows found.")
        return nil
    }

    fmt.Printf("%-40s %-10s %-20s %s\n", "ID", "Status", "Created At", "Name")
    for _, workflow := range workflows {
        fmt.Printf("%-40s %-10s %-20s %s\n", workflow.ID.String(), workflow.Status, workflow.CreatedAt.Format("2006-01-02 15:04:05"), workflow.Name)
    }
    return nil
}

func runWorkflowShow(cmd *cobra.Command, tenantID *string, workflowID string) error {
    tenantUUID, err := requireTenantID(*tenantID)
    if err != nil {
        return err
    }

    svc, err := loadServices()
    if err != nil {
        return err
    }

    id, err := uuid.Parse(workflowID)
    if err != nil {
        return fmt.Errorf("invalid workflow ID: %w", err)
    }

    workflow, err := svc.Workflow.GetWorkflow(tenantUUID, id)
    if err != nil {
        return fmt.Errorf("failed to retrieve workflow: %w", err)
    }

    output, err := json.MarshalIndent(workflow, "", "  ")
    if err != nil {
        return err
    }
    fmt.Println(string(output))
    return nil
}

func runWorkflowDelete(cmd *cobra.Command, tenantID *string, workflowID string) error {
    tenantUUID, err := requireTenantID(*tenantID)
    if err != nil {
        return err
    }

    svc, err := loadServices()
    if err != nil {
        return err
    }

    id, err := uuid.Parse(workflowID)
    if err != nil {
        return fmt.Errorf("invalid workflow ID: %w", err)
    }

    if err := svc.Workflow.DeleteWorkflow(tenantUUID, id); err != nil {
        return fmt.Errorf("failed to delete workflow: %w", err)
    }
    fmt.Println("Workflow deleted")
    return nil
}

func runWorkflowRun(cmd *cobra.Command, tenantID *string, workflowID, input string, async bool) error {
    tenantUUID, err := requireTenantID(*tenantID)
    if err != nil {
        return err
    }

    svc, err := loadServices()
    if err != nil {
        return err
    }

    id, err := uuid.Parse(workflowID)
    if err != nil {
        return fmt.Errorf("invalid workflow ID: %w", err)
    }

    workflow, err := svc.Workflow.GetWorkflow(tenantUUID, id)
    if err != nil {
        return fmt.Errorf("failed to retrieve workflow: %w", err)
    }

    if async {
        if err := svc.Workflow.EnqueueWorkflowRun(tenantUUID, id, input); err != nil {
            return fmt.Errorf("failed to enqueue workflow run: %w", err)
        }
        fmt.Println("Workflow execution enqueued")
        return nil
    }

    result, err := svc.Workflow.ExecuteWorkflow(context.Background(), tenantUUID, workflow, input)
    if err != nil {
        return fmt.Errorf("workflow execution failed: %w", err)
    }
    fmt.Println("Workflow output:")
    fmt.Println(result)
    return nil
}

func loadServices() (*services.Services, error) {
    cfg, err := config.LoadConfig(".")
    if err != nil {
        return nil, err
    }
    db, err := database.Connect(cfg)
    if err != nil {
        return nil, err
    }
    repo := repositories.NewRepository(db)
    logger, _ := zap.NewProduction()
    return services.NewServices(cfg, repo, logger), nil
}

func requireTenantID(tenantID string) (uuid.UUID, error) {
    if tenantID == "" {
        tenantID = os.Getenv("TENANT_ID")
    }
    if tenantID == "" {
        return uuid.Nil, errors.New("tenant ID is required; pass --tenant or set TENANT_ID")
    }
    return uuid.Parse(tenantID)
}

func buildActionsFromAgents(tenantID uuid.UUID, agentService *services.AgentService, identifiers, prompt string) ([]workflowAction, error) {
    ids := strings.Split(identifiers, ",")
    var actions []workflowAction
    agents, err := agentService.ListAgents(tenantID)
    if err != nil {
        return nil, fmt.Errorf("failed to list agents: %w", err)
    }

    for _, identifier := range ids {
        identifier = strings.TrimSpace(identifier)
        if identifier == "" {
            continue
        }

        var agent *models.Agent
        if parsed, err := uuid.Parse(identifier); err == nil {
            agent, err = agentService.GetAgent(tenantID, parsed)
            if err != nil {
                return nil, fmt.Errorf("unable to find agent by ID %s: %w", identifier, err)
            }
        } else {
            for i := range agents {
                if strings.EqualFold(agents[i].Name, identifier) || strings.EqualFold(agents[i].Type, identifier) {
                    agent = &agents[i]
                    break
                }
            }
            if agent == nil {
                return nil, fmt.Errorf("unable to find agent by name: %s", identifier)
            }
        }

        actions = append(actions, workflowAction{
            Name:    agent.Name,
            AgentID: agent.ID.String(),
            Prompt:  prompt,
        })
    }
    if len(actions) == 0 {
        return nil, errors.New("no workflow actions were created from agents")
    }
    return actions, nil
}
