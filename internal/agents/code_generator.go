
package agents

import (
	"context"
	"fmt"

	"velora/internal/services"
)

// CodeGenerator is an agent that generates code.
type CodeGenerator struct {
	llm services.LLMService
}

// NewCodeGenerator creates a new code generator agent.
func NewCodeGenerator(llm services.LLMService) *CodeGenerator {
	return &CodeGenerator{llm: llm}
}

// Name returns the name of the agent.
func (a *CodeGenerator) Name() string {
	return "CodeGenerator"
}

// Description returns the description of the agent.
func (a *CodeGenerator) Description() string {
	return "Generates code based on a prompt."
}

// Run executes the agent.
func (a *CodeGenerator) Run(ctx context.Context, input string) (string, error) {
	prompt := fmt.Sprintf("Generate Go code for the following task: %s", input)
	return a.llm.Generate(prompt, "gemini-1.5-pro", 0.7, 2048)
}
