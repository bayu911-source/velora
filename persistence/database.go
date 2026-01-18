
package persistence

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	// DB is the global database connection.
	DB *sql.DB
)

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
