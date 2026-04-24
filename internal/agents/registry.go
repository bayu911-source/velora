package agents

import (
	"fmt"

	"velora/internal/services"
)

// Registry stores the available agents.
type Registry struct {
	agents map[string]Agent
	llm    *services.LLM
}

// NewRegistry creates a new agent registry.
func NewRegistry(llm *services.LLM) *Registry {
	return &Registry{
		agents: make(map[string]Agent),
		llm:    llm,
	}
}

// Register registers a new agent.
func (r *Registry) Register(agent Agent) {
	r.agents[agent.Name()] = agent
}

// Get returns an agent by name.
func (r *Registry) Get(name string) (Agent, error) {
	agent, ok := r.agents[name]
	if !ok {
		return nil, fmt.Errorf("agent %s not found", name)
	}
	return agent, nil
}

// List returns a list of all available agents.
func (r *Registry) List() []Agent {
	var agents []Agent
	for _, agent := range r.agents {
		agents = append(agents, agent)
	}
	return agents
}

// LLM returns the LLM service.
func (r *Registry) LLM() *services.LLM {
	return r.llm
}
