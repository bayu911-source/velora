package workflow

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// State defines the possible states of a workflow.
type State string

const (
	// StatePending means the workflow has not started yet.
	StatePending State = "PENDING"
	// StateRunning means the workflow is in progress.
	StateRunning State = "RUNNING"
	// StateCompleted means the workflow has finished successfully.
	StateCompleted State = "COMPLETED"
	// StateFailed means the workflow has failed.
	StateFailed State = "FAILED"
)

// Workflow defines a sequence of agents to be executed.
// This struct is now primarily a data container.
type Workflow struct {
	ID          string
	Name        string
	Description string
	State       State
	CreatedAt   time.Time
	UpdatedAt   time.Time
	steps       []*Step
}

// Step represents a single step in a workflow.
type Step struct {
	ID         int
	WorkflowID string
	AgentName  string
	Input      string
	Output     string
	ExecutedAt sql.NullTime
}

// New creates a new Workflow instance in memory.
// It no longer saves the workflow to the database.
func New(name, description string, agentNames []string) *Workflow {
	w := &Workflow{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		State:       StatePending,
	}

	for _, agentName := range agentNames {
		step := &Step{
			WorkflowID: w.ID,
			AgentName:  agentName,
		}
		w.steps = append(w.steps, step)
	}

	return w
}

// Steps returns the steps of the workflow.
func (w *Workflow) Steps() []*Step {
    return w.steps
}

// SetSteps sets the steps of the workflow.
func (w *Workflow) SetSteps(steps []*Step) {
	w.steps = steps
}
