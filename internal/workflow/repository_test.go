package workflow

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func createTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	createWorkflowsTableSQL := `
		CREATE TABLE IF NOT EXISTS workflows (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			state TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`
	if _, err := db.Exec(createWorkflowsTableSQL); err != nil {
		t.Fatalf("failed to create workflows table: %v", err)
	}

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
	if _, err := db.Exec(createWorkflowStepsTableSQL); err != nil {
		t.Fatalf("failed to create workflow_steps table: %v", err)
	}

	return db
}

func TestRepository_SaveAndDelete(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	workflow := New("test-workflow", "", []string{"agent1"})
	workflow.Steps()[0].Input = "initial"
	workflow.Steps()[0].Output = "first-output"

	if err := repo.Save(workflow); err != nil {
		t.Fatalf("expected save to succeed, got %v", err)
	}

	workflow.Steps()[0].Output = "updated-output"
	if err := repo.Save(workflow); err != nil {
		t.Fatalf("expected save update to succeed, got %v", err)
	}

	loaded, err := repo.FindByID(workflow.ID)
	if err != nil {
		t.Fatalf("failed to load workflow: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected workflow to be found after save")
	}
	if len(loaded.Steps()) != 1 {
		t.Fatalf("expected 1 step after update, got %d", len(loaded.Steps()))
	}
	if loaded.Steps()[0].Output != "updated-output" {
		t.Fatalf("expected updated output %q, got %q", "updated-output", loaded.Steps()[0].Output)
	}

	if err := repo.DeleteByID(workflow.ID); err != nil {
		t.Fatalf("expected delete to succeed, got %v", err)
	}

	deleted, err := repo.FindByID(workflow.ID)
	if err != nil {
		t.Fatalf("error loading workflow after delete: %v", err)
	}
	if deleted != nil {
		t.Fatal("expected workflow to be deleted")
	}
}
