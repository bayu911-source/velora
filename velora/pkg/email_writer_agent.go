
package pkg

import (
	"fmt"
)

type EmailWriterAgent struct {
	gemini *GeminiService
}

func NewEmailWriterAgent(gemini *GeminiService) *EmailWriterAgent {
	return &EmailWriterAgent{
		gemini: gemini,
	}
}

func (a *EmailWriterAgent) Name() string {
	return "email_writer"
}

// Run generates an email. It uses the input (analysis from the previous step)
// and the initial topic from memory to create a comprehensive prompt.
func (a *EmailWriterAgent) Run(memory *MemoryManager, analysisResult string) (string, error) {
	// Retrieve the original input that started the pipeline.
	// The runner stores this under the key "initial_input".
	initialTopic := "the user's request" // Fallback text
	if topic, ok := memory.Get("initial_input"); ok {
		if topicStr, ok := topic.(string); ok {
			initialTopic = topicStr
		}
	}

	// Combine the initial topic with the analysis from the previous step.
	prompt := fmt.Sprintf(
		"Based on the initial topic of '%s' and the following analysis:\n\n---\n%s\n---\n\nDraft a professional email.",
		initialTopic,
		analysisResult,
	)

	return a.gemini.Generate(prompt, 0.7, 1024)
}
