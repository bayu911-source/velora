
package agents

import (
	"context"
	"fmt"

	"google.golang.org/genai"
	"velora/internal/services"
)

// CodeGenerator is an agent that generates code.
type CodeGenerator struct{}

// NewCodeGenerator creates a new code generator agent.
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
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
func (a *CodeGenerator) Run(ctx context.Context, llm *services.LLM, input string) (string, error) {
	prompt := fmt.Sprintf("Generate Go code for the following task: %s", input)
	resp, err := llm.Model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate code: %w", err)
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(txt), nil
		}
	}

	return "", fmt.Errorf("no code response found")
}
