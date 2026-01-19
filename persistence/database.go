
package persistence

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// DB is the global database connection.
	DB *sql.DB
)

// Workflow represents the state of a workflow.
type Workflow struct {
	ID        string
	Name      string
	State     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// WorkflowStep represents a single step in a workflow.
type WorkflowStep struct {
	ID         int64
	WorkflowID string
	AgentName  string
	Input      string
	Output     string
	ExecutedAt time.Time
}

// InitDB initializes the database connection.
func InitDB() {
	// Get the database path from the environment variable or use a default value.
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "velora.db"
	}

	// Open the database connection.
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// Create the necessary tables if they don't exist.
	createTables()
}

// createTables creates the database tables if they don't already exist.
func createTables() {
	// SQL statement to create the workflows table.
	createWorkflowsTableSQL := `
		CREATE TABLE IF NOT EXISTS workflows (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			state TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`

	// Execute the SQL statement.
	_, err := DB.Exec(createWorkflowsTableSQL)
	if err != nil {
		log.Fatalf("failed to create workflows table: %v", err)
	}

	// SQL statement to create the workflow_steps table.
	createWorkflowStepsTableSQL := `
		CREATE TABLE IF NOT EXISTS workflow_steps (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			workflow_id TEXT NOT NULL,
			agent_name TEXT NOT NULL,
			input TEXT,
			output TEXT,
			executed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(workflow_id) REFERENCES workflows(id)
		);
	`

	// Execute the SQL statement.
	_, err = DB.Exec(createWorkflowStepsTableSQL)
	if err != nil {
		log.Fatalf("failed to create workflow_steps table: %v", err)
	}
}

// CreateWorkflow creates a new workflow in the database.
func CreateWorkflow(name string, initialState string) (*Workflow, error) {
	workflow := &Workflow{
		ID:    uuid.New().String(),
		Name:  name,
		State: initialState,
	}

	stmt, err := DB.Prepare("INSERT INTO workflows (id, name, state) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(workflow.ID, workflow.Name, workflow.State)
	if err != nil {
		return nil, err
	}

	return workflow, nil
}

// GetWorkflow retrieves a workflow from the database by its ID.
func GetWorkflow(id string) (*Workflow, error) {
	row := DB.QueryRow("SELECT id, name, state, created_at, updated_at FROM workflows WHERE id = ?", id)

	workflow := &Workflow{}
	err := row.Scan(&workflow.ID, &workflow.Name, &workflow.State, &workflow.CreatedAt, &workflow.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	return workflow, nil
}

// UpdateWorkflowState updates the state of a workflow.
func UpdateWorkflowState(id, state string) error {
	stmt, err := DB.Prepare("UPDATE workflows SET state = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(state, time.Now(), id)
	return err
}

// CreateWorkflowStep creates a new workflow step in the database.
func CreateWorkflowStep(workflowID, agentName, input, output string) error {
	stmt, err := DB.Prepare("INSERT INTO workflow_steps (workflow_id, agent_name, input, output) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(workflowID, agentName, input, output)
	return err
}
