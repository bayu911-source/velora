package server

import (
	"encoding/json"
	"net/http"

	"velora/internal/agents"
	"velora/internal/workflows"
)

// NewServer membuat instance server http baru.
func NewServer() *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/agents", cors(agentsHandler))
	mux.HandleFunc("/workflows", cors(workflowsHandler))

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func agentsHandler(w http.ResponseWriter, r *http.Request) {
	agentList := agents.List()
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

func workflowsHandler(w http.ResponseWriter, r *http.Request) {
	workflowList := workflows.List()
	response := []map[string]string{}

	for _, workflow := range workflowList {
		response = append(response, map[string]string{
			"name":        workflow.Name,
			"description": workflow.Description,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
