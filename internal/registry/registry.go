
package registry

import (
	"log"

	"velora/internal/agents"
	"velora/internal/services"
)

var AgentRegistry *agents.Registry

func InitAgents(llm services.LLM) {
	AgentRegistry = agents.NewRegistry()

	// Register all agents
	if err := AgentRegistry.Register(agents.NewTextAnalysisAgent(llm)); err != nil {
		log.Fatalf("Gagal mendaftarkan agen: %v", err)
	}
	if err := AgentRegistry.Register(agents.NewEmailWriterAgent(llm)); err != nil {
		log.Fatalf("Gagal mendaftarkan agen: %v", err)
	}
	if err := AgentRegistry.Register(agents.NewDataExtractorAgent(llm)); err != nil {
		log.Fatalf("Gagal mendaftarkan agen: %v", err)
	}
	if err := AgentRegistry.Register(agents.NewWebResearchAgent(llm)); err != nil {
		log.Fatalf("Gagal mendaftarkan agen: %v", err)
	}
	if err := AgentRegistry.Register(agents.NewAppBuilderAgent(llm)); err != nil {
		log.Fatalf("Gagal mendaftarkan agen: %v", err)
	}
	if err := AgentRegistry.Register(agents.NewMultimodalAgent(llm)); err != nil {
		log.Fatalf("Gagal mendaftarkan agen: %v", err)
	}
	if err := AgentRegistry.Register(agents.NewChatAgent(llm)); err != nil {
		log.Fatalf("Gagal mendaftarkan agen: %v", err)
	}
}
