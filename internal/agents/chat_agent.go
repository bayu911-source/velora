
package agents

import (
	"context"
	"fmt"

	"velora/internal/services"
)

// ChatAgent is an agent that can have a multi-step conversation.
type ChatAgent struct {
	llm services.LLMService
}

// NewChatAgent creates a new ChatAgent.
func NewChatAgent(llm services.LLMService) *ChatAgent {
	return &ChatAgent{llm: llm}
}

// Name returns the name of the agent.
func (a *ChatAgent) Name() string {
	return "chat"
}

// Description returns the description of the agent.
func (a *ChatAgent) Description() string {
	return "Starts an interactive chat session with the AI."
}

// Run executes the agent.
func (a *ChatAgent) Run(ctx context.Context, input string) (string, error) {
	// For a chat agent, we want to stream the response.
	// We'll use the GenerateStream method of the LLM service.
	stream, errChan := a.llm.GenerateStream(input, "gemini-1.5-pro", 0.7, 2048)

	var response string
	for {
		select {
		case res, ok := <-stream:
			if !ok {
				return response, nil
			}
			response += res
			fmt.Print(res)
		case err := <-errChan:
			return "", fmt.Errorf("failed to generate response: %w", err)
		}
	}
}
