
package agents

import (
	"fmt"

	"velora/internal/services"
)

// Registry is a registry for all available agents.
type Registry struct {
	agents map[string]Agent
}

// NewRegistry creates a new Registry.
func NewRegistry(llm services.LLM) *Registry {
	r := &Registry{agents: make(map[string]Agent)}

	// Register all agents
	r.Register(NewTextAnalysisAgent(llm))
	r.Register(NewEmailWriterAgent(llm))
	r.Register(NewDataExtractorAgent(llm))
	r.Register(NewWebResearchAgent(llm))
	r.Register(NewAppBuilderAgent(llm))
	r.Register(NewChatAgent(llm))
	r.Register(NewMultimodalAgent(llm))

	return r
}

// Register registers a new agent.
func (r *Registry) Register(agent Agent) error {
	if _, exists := r.agents[agent.Name()]; exists {
		return fmt.Errorf("agent with name '%s' already registered", agent.Name())
	}
	r.agents[agent.Name()] = agent
	return nil
}

// Get returns the agent with the given name.
func (r *Registry) Get(name string) (Agent, error) {
	agent, exists := r.agents[name]
	if !exists {
		return nil, fmt.Errorf("agent with name '%s' not found", name)
	}
	return agent, nil
}

// List returns a list of all registered agents.
func (r *Registry) List() []Agent {
	var agents []Agent
	for _, agent := range r.agents {
		agents = append(agents, agent)
	}
	return agents
}
