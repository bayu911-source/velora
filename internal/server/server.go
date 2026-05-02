package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"velora/internal/agents"
	"velora/internal/workflow"
	"velora/persistence"
)

// NewServer creates a new HTTP server for Velora.
func NewServer(registry *agents.Registry) *http.Server {
	repo := workflow.NewRepository(persistence.DB)
	engine := workflow.NewEngine(registry, repo)
	mux := http.NewServeMux()

	mux.HandleFunc("/api/agents", cors(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		agentsHandler(w, r, registry)
	}))

	mux.HandleFunc("/api/workflows", cors(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			workflowsHandler(w, repo)
		case http.MethodPost:
			createWorkflowHandler(w, r, repo)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/workflows/", cors(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/workflows/")
		parts := strings.Split(path, "/")
		if len(parts) == 0 || parts[0] == "" {
			http.NotFound(w, r)
			return
		}

		workflowID := parts[0]
		switch {
		case len(parts) == 1 && r.Method == http.MethodGet:
			workflowDetailHandler(w, r, repo, workflowID)
			return
		case len(parts) == 1 && r.Method == http.MethodDelete:
			deleteWorkflowHandler(w, repo, workflowID)
			return
		case len(parts) == 2 && parts[1] == "run" && r.Method == http.MethodPost:
			workflowRunHandler(w, r, repo, engine, workflowID)
			return
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))

	staticDir := "ui/build"
	if _, err := os.Stat(filepath.Join(staticDir, "index.html")); err == nil {
		fileServer := http.FileServer(http.Dir(staticDir))
		mux.Handle("/static/", fileServer)
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.NotFound(w, r)
				return
			}

			requestedPath := filepath.Join(staticDir, filepath.Clean(r.URL.Path))
			if r.URL.Path != "/" {
				if _, err := os.Stat(requestedPath); err == nil {
					http.ServeFile(w, r, requestedPath)
					return
				}
			}

			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
		})
	}

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func agentsHandler(w http.ResponseWriter, r *http.Request, registry *agents.Registry) {
	agentList := registry.List()
	response := []map[string]string{}

	for _, agent := range agentList {
		response = append(response, map[string]string{
			"name":        agent.Name(),
			"description": agent.Description(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func workflowsHandler(w http.ResponseWriter, repo *workflow.Repository) {
	workflows, err := repo.ListAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list workflows: %v", err), http.StatusInternalServerError)
		return
	}

	response := make([]map[string]interface{}, 0, len(workflows))
	for _, wfl := range workflows {
		response = append(response, map[string]interface{}{
			"id":          wfl.ID,
			"name":        wfl.Name,
			"description": wfl.Description,
			"state":       wfl.State,
			"created_at":  wfl.CreatedAt,
			"updated_at":  wfl.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createWorkflowHandler(w http.ResponseWriter, r *http.Request, repo *workflow.Repository) {
	var body struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Agents      []string `json:"agents"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if body.Name == "" || len(body.Agents) == 0 {
		http.Error(w, "name and agents are required", http.StatusBadRequest)
		return
	}

	newWorkflow := workflow.New(body.Name, body.Description, body.Agents)
	if err := repo.Save(newWorkflow); err != nil {
		http.Error(w, fmt.Sprintf("failed to save workflow: %v", err), http.StatusInternalServerError)
		return
	}

	wResponse := map[string]interface{}{
		"id":          newWorkflow.ID,
		"name":        newWorkflow.Name,
		"description": newWorkflow.Description,
		"state":       newWorkflow.State,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wResponse)
}

func workflowDetailHandler(w http.ResponseWriter, r *http.Request, repo *workflow.Repository, id string) {
	workflow, err := repo.FindByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to find workflow: %v", err), http.StatusInternalServerError)
		return
	}
	if workflow == nil {
		http.NotFound(w, r)
		return
	}

	response := map[string]interface{}{
		"id":          workflow.ID,
		"name":        workflow.Name,
		"description": workflow.Description,
		"state":       workflow.State,
		"created_at":  workflow.CreatedAt,
		"updated_at":  workflow.UpdatedAt,
		"steps":       workflow.Steps(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deleteWorkflowHandler(w http.ResponseWriter, repo *workflow.Repository, id string) {
	if err := repo.DeleteByID(id); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete workflow: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func workflowRunHandler(w http.ResponseWriter, r *http.Request, repo *workflow.Repository, engine *workflow.Engine, id string) {
	workflowObj, err := repo.FindByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to find workflow: %v", err), http.StatusInternalServerError)
		return
	}
	if workflowObj == nil {
		http.NotFound(w, r)
		return
	}

	var body struct {
		Input string `json:"input"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	output, err := engine.Run(context.Background(), workflowObj, body.Input)
	if err != nil {
		http.Error(w, fmt.Sprintf("workflow execution failed: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":     workflowObj.ID,
		"output": output,
		"state":  workflowObj.State,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
