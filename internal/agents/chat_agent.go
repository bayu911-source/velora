package agents

import (
	"context"

	"velora/internal/services"
)

// ChatAgent is a conversational agent.
type ChatAgent struct {
	llm *services.LLM
}

// NewChatAgent creates a new ChatAgent.
func NewChatAgent(llm *services.LLM) *ChatAgent {
	return &ChatAgent{llm: llm}
}

// Name returns the name of the agent.
func (a *ChatAgent) Name() string {
	return "ChatAgent"
}

// Description returns the description of the agent.
func (a *ChatAgent) Description() string {
	return "A conversational agent that can answer questions."
}

// Execute runs the agent.
func (a *ChatAgent) Execute(ctx context.Context, input string) (string, error) {
	return a.llm.Generate(ctx, input)
}
