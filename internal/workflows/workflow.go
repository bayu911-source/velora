
package workflows

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"velora/internal/agents"
	"velora/persistence"

	"github.com/google/uuid"
)

// WorkflowState mendefinisikan kemungkinan status alur kerja.
type WorkflowState string

const (
	// WorkflowStatePending berarti alur kerja belum dimulai.
	WorkflowStatePending WorkflowState = "PENDING"
	// WorkflowStateRunning berarti alur kerja sedang dalam proses.
	WorkflowStateRunning WorkflowState = "RUNNING"
	// WorkflowStateCompleted berarti alur kerja telah berhasil diselesaikan.
	WorkflowStateCompleted WorkflowState = "COMPLETED"
	// WorkflowStateFailed berarti alur kerja mengalami kesalahan.
	WorkflowStateFailed WorkflowState = "FAILED"
)

// Workflow mendefinisikan urutan agen untuk dieksekusi.
type Workflow struct {
	ID          string
	Name        string
	State       WorkflowState
	CreatedAt   time.Time
	UpdatedAt   time.Time
	steps       []*WorkflowStep
	currentStep int
}

// WorkflowStep mewakili satu langkah dalam alur kerja.
type WorkflowStep struct {
	ID         int
	WorkflowID string
	AgentName  string
	Input      string
	Output     string
	ExecutedAt time.Time
}

// New membuat instance Alur Kerja baru.
func New(name string, agentNames []string) (*Workflow, error) {
	w := &Workflow{
		ID:    uuid.New().String(),
		Name:  name,
		State: WorkflowStatePending,
	}

	for _, agentName := range agentNames {
		step := &WorkflowStep{
			WorkflowID: w.ID,
			AgentName:  agentName,
		}
		w.steps = append(w.steps, step)
	}

	if err := w.Save(); err != nil {
		return nil, err
	}

	return w, nil
}

// Run mengeksekusi alur kerja.
func (w *Workflow) Run(ctx context.Context, agentRegistry *agents.Registry, initialInput string) (string, error) {
	w.State = WorkflowStateRunning
	if err := w.Save(); err != nil {
		return "", fmt.Errorf("gagal menyimpan status alur kerja awal: %w", err)
	}

	var currentInput = initialInput
	var output string

	steps, err := w.GetSteps()
	if err != nil {
		return "", fmt.Errorf("gagal mengambil langkah alur kerja: %w", err)
	}

	for _, step := range steps {
		agent, err := agentRegistry.Get(step.AgentName)
		if err != nil {
			return "", fmt.Errorf("gagal mengambil agen '%s': %w", step.AgentName, err)
		}
		fmt.Printf("Menjalankan agen: %s\n", agent.Name())

		step.Input = currentInput
		step.ExecutedAt = time.Now()

		output, err = agent.Run(ctx, currentInput)
		if err != nil {
			w.State = WorkflowStateFailed
			step.Output = err.Error()
			if err := w.SaveStep(step); err != nil {
				return "", fmt.Errorf("gagal menyimpan langkah alur kerja yang gagal: %w", err)
			}
			if err := w.Save(); err != nil {
				return "", fmt.Errorf("gagal menyimpan status alur kerja yang gagal: %w", err)
			}
			return "", fmt.Errorf("gagal menjalankan agen %s: %w", agent.Name(), err)
		}
		currentInput = output
		step.Output = output

		if err := w.SaveStep(step); err != nil {
			return "", fmt.Errorf("gagal menyimpan langkah alur kerja yang berhasil: %w", err)
		}
	}

	w.State = WorkflowStateCompleted
	if err := w.Save(); err != nil {
		return "", fmt.Errorf("gagal menyimpan status alur kerja yang telah selesai: %w", err)
	}

	return output, nil
}

// Save menyimpan status alur kerja ke database.
func (w *Workflow) Save() error {
	query := `
		INSERT INTO workflows (id, name, state, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			state = excluded.state,
			updated_at = excluded.updated_at;
	`

	w.UpdatedAt = time.Now()
	if w.CreatedAt.IsZero() {
		w.CreatedAt = w.UpdatedAt
	}

	_, err := persistence.DB.Exec(query, w.ID, w.Name, w.State, w.CreatedAt, w.UpdatedAt)
	return err
}

// SaveStep menyimpan langkah alur kerja tunggal ke database.
func (w *Workflow) SaveStep(step *WorkflowStep) error {
	query := `
		INSERT INTO workflow_steps (workflow_id, agent_name, input, output, executed_at)
		VALUES (?, ?, ?, ?, ?);
	`
	_, err := persistence.DB.Exec(query, w.ID, step.AgentName, step.Input, step.Output, step.ExecutedAt)
	return err
}

// GetWorkflowByID mengambil alur kerja dari database berdasarkan ID-nya.
func GetWorkflowByID(id string) (*Workflow, error) {
	query := `SELECT id, name, state, created_at, updated_at FROM workflows WHERE id = ?`
	row := persistence.DB.QueryRow(query, id)

	w := &Workflow{}
	err := row.Scan(&w.ID, &w.Name, &w.State, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No workflow found with that ID
		}
		return nil, err
	}

	return w, nil
}

// GetSteps mengambil langkah-langkah untuk alur kerja tertentu dari database.
func (w *Workflow) GetSteps() ([]*WorkflowStep, error) {
	query := `SELECT id, agent_name, input, output, executed_at FROM workflow_steps WHERE workflow_id = ? ORDER BY id ASC`
	rows, err := persistence.DB.Query(query, w.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []*WorkflowStep
	for rows.Next() {
		step := &WorkflowStep{WorkflowID: w.ID}
		if err := rows.Scan(&step.ID, &step.AgentName, &step.Input, &step.Output, &step.ExecutedAt); err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}

	return steps, nil
}

// ListWorkflows mengambil daftar semua alur kerja dari database.
func ListWorkflows() ([]*Workflow, error) {
	query := `SELECT id, name, state, created_at, updated_at FROM workflows`
	rows, err := persistence.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []*Workflow
	for rows.Next() {
		w := &Workflow{}
		if err := rows.Scan(&w.ID, &w.Name, &w.State, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		workflows = append(workflows, w)
	}

	return workflows, nil
}
