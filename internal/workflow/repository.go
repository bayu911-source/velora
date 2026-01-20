package workflow

import (
	"database/sql"
	"time"

	"velora/persistence"
)

// Repository handles the database operations for workflows.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new workflow repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Save saves the workflow and its steps to the database.
func (r *Repository) Save(w *Workflow) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

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

	if _, err := tx.Exec(query, w.ID, w.Name, w.State, w.CreatedAt, w.UpdatedAt); err != nil {
		tx.Rollback()
		return err
	}

	stepQuery := `
		INSERT INTO workflow_steps (workflow_id, agent_name, input, output, executed_at)
		VALUES (?, ?, ?, ?, ?);
	`
	for _, step := range w.Steps() {
		if _, err := tx.Exec(stepQuery, w.ID, step.AgentName, step.Input, step.Output, step.ExecutedAt); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// FindByID retrieves a workflow by its ID.
func (r *Repository) FindByID(id string) (*Workflow, error) {
	query := `SELECT id, name, state, created_at, updated_at FROM workflows WHERE id = ?`
	row := r.db.QueryRow(query, id)

	w := &Workflow{}
	err := row.Scan(&w.ID, &w.Name, &w.State, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}

	steps, err := r.findStepsByWorkflowID(w.ID)
	if err != nil {
		return nil, err
	}
	w.SetSteps(steps)

	return w, nil
}

// findStepsByWorkflowID retrieves the steps for a given workflow ID.
func (r *Repository) findStepsByWorkflowID(workflowID string) ([]*Step, error) {
	query := `SELECT id, agent_name, input, output, executed_at FROM workflow_steps WHERE workflow_id = ? ORDER BY id ASC`
	rows, err := r.db.Query(query, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []*Step
	for rows.Next() {
		step := &Step{WorkflowID: workflowID}
		if err := rows.Scan(&step.ID, &step.AgentName, &step.Input, &step.Output, &step.ExecutedAt); err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}

	return steps, nil
}

// ListAll retrieves all workflows from the database.
func (r *Repository) ListAll() ([]*Workflow, error) {
	query := `SELECT id, name, state, created_at, updated_at FROM workflows`
	rows, err := r.db.Query(query)
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
