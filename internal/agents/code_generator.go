package agents

import (
	"context"

	"velora/internal/services"
)

// CodeGenerator is an agent that generates code.
type CodeGenerator struct {
	llm *services.LLM
}

// NewCodeGenerator creates a new CodeGenerator.
func NewCodeGenerator(llm *services.LLM) *CodeGenerator {
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

// Execute runs the agent.
func (a *CodeGenerator) Execute(ctx context.Context, input string) (string, error) {
	return a.llm.Generate(ctx, "Generate code for: "+input)
}
