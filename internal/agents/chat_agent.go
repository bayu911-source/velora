package agents

import (
	"context"
)

// ChatAgent is a conversational agent.
type ChatAgent struct {
	llm LLMService
}

// NewChatAgent creates a new ChatAgent.
func NewChatAgent(llm LLMService) *ChatAgent {
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
