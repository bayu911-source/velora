
package main

import (
	"velora/cmd"
)

func main() {
	cmd.Execute()
	// Initialize Gemini Service
	// gemini, err := services.NewGeminiService("YOUR_GEMINI_API_KEY")
	// if err != nil {
	// 	log.Fatalf("Failed to create Gemini service: %v", err)
	// }

	// // Initialize Memory Manager
	// memoryManager := memory.NewMemoryManager("memory.json")

	// // Register Agents
	// agentRegistry := make(map[string]agents.Agent)
	// agentRegistry["text_analysis"] = agents.NewTextAnalysisAgent(gemini)
	// agentRegistry["email_writer"] = agents.NewEmailWriterAgent(gemini)
	// agentRegistry["data_extractor"] = agents.NewDataExtractorAgent(gemini)
	// agentRegistry["web_research"] = agents.NewWebResearchAgent(gemini)

	// // Initialize Pipeline Runner
	// pipelineRunner := workflow.NewPipelineRunner(agentRegistry)

	// // HTTP Handlers
	// http.HandleFunc("/agent/run", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != http.MethodPost {
	// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	// 		return
	// 	}

	// 	body, err := ioutil.ReadAll(r.Body)
	// 	if err != nil {
	// 		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	var request struct {
	// 		Agent string `json:"agent"`
	// 		Input string `json:"input"`
	// 	}

	// 	if err := json.Unmarshal(body, &request); err != nil {
	// 		http.Error(w, "Failed to unmarshal request", http.StatusBadRequest)
	// 		return
	// 	}

	// 	agent, ok := agentRegistry[request.Agent]
	// 	if !ok {
	// 		http.Error(w, fmt.Sprintf("Agent %q not found", request.Agent), http.StatusBadRequest)
	// 		return
	// 	}

	// 	output, err := agent.Run(request.Input)
	// 	if err != nil {
	// 		http.Error(w, fmt.Sprintf("Agent failed: %v", err), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	response := map[string]string{"output": output}
	// 	json.NewEncoder(w).Encode(response)
	// })

	// http.HandleFunc("/workflow/run", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != http.MethodPost {
	// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	// 		return
	// 	}

	// 	body, err := ioutil.ReadAll(r.Body)
	// 	if err != nil {
	// 		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	outputs, err := pipelineRunner.Run(string(body))
	// 	if err != nil {
	// 		http.Error(w, fmt.Sprintf("Workflow failed: %v", err), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	json.NewEncoder(w).Encode(outputs)
	// })

	// http.HandleFunc("/memory/save", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != http.MethodPost {
	// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	// 		return
	// 	}

	// 	body, err := ioutil.ReadAll(r.Body)
	// 	if err != nil {
	// 		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	var request struct {
	// 		Key   string `json:"key"`
	// 		Value string `json:"value"`
	// 	}

	// 	if err := json.Unmarshal(body, &request); err != nil {
	// 		http.Error(w, "Failed to unmarshal request", http.StatusBadRequest)
	// 		return
	// 	}

	// 	if err := memoryManager.SaveMemory(request.Key, request.Value); err != nil {
	// 		http.Error(w, fmt.Sprintf("Failed to save memory: %v", err), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.WriteHeader(http.StatusOK)
	// })

	// http.HandleFunc("/memory/get", func(w http.ResponseWriter, r *http.Request) {
	// 	key := r.URL.Query().Get("key")
	// 	value, ok := memoryManager.GetMemory(key)
	// 	if !ok {
	// 		http.Error(w, "Key not found in memory", http.StatusNotFound)
	// 		return
	// 	}

	// 	response := map[string]string{"value": value}
	// 	json.NewEncoder(w).Encode(response)
	// })

	// // Start Server
	// log.Println("Server starting on port 8080...")
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	log.Fatalf("Server failed: %v", err)
	// }
}
