
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/velora-id/velora/internal/workflow"
	"github.com/velora-id/velora/pkg"
)

func main() {
	// 1. Initialize Gemini Service
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set")
	}
	gemini, err := pkg.NewGeminiService(geminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to create Gemini service: %v", err)
	}

	// 2. Create and Configure the Runner
	runner := workflow.NewRunner()

	// 3. Register all agents with the runner
	runner.RegisterAgent(pkg.NewWebResearchAgent(gemini))
	runner.RegisterAgent(pkg.NewTextAnalysisAgent(gemini))
	runner.RegisterAgent(pkg.NewEmailWriterAgent(gemini))

	// 4. Define the HTTP handler for running pipelines
	http.HandleFunc("/run/", func(w http.ResponseWriter, r *http.Request) {
		pipelineName := strings.TrimPrefix(r.URL.Path, "/run/")
		if pipelineName == "" {
			http.Error(w, "Pipeline name not specified", http.StatusBadRequest)
			return
		}

		// Read the pipeline definition from the corresponding JSON file
		jsonFile, err := os.Open(pipelineName + ".json")
		if err != nil {
			http.Error(w, fmt.Sprintf("Pipeline '%s' not found", pipelineName), http.StatusNotFound)
			return
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		var pipeline workflow.Pipeline
		json.Unmarshal(byteValue, &pipeline)

		// Execute the pipeline
		result, err := runner.Run(&pipeline)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error executing pipeline: %v", err), http.StatusInternalServerError)
			return
		}

		// Return the result
		fmt.Fprint(w, result)
	})

	// 5. Start the HTTP server
	addr := ":8080"
	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
