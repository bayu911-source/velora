
package registry

import (
	"velora/internal/agents"
	"velora/internal/services"
)

var AgentRegistry *agents.Registry

func InitAgents(llm *services.LLM) {
	AgentRegistry = agents.NewRegistry(llm)

	// Register all agents
	AgentRegistry.Register(agents.NewTextAnalysisAgent(llm))
	AgentRegistry.Register(agents.NewEmailWriterAgent(llm))
	AgentRegistry.Register(agents.NewDataExtractorAgent(llm))
	AgentRegistry.Register(agents.NewWebResearchAgent(llm))
	AgentRegistry.Register(agents.NewAppBuilderAgent(llm))
	AgentRegistry.Register(agents.NewChatAgent(llm))
}
