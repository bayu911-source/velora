
package pkg

import (
	"context"
	"os"

	"google.golang.org/genai"
)

type GeminiService struct {
	client *genai.Client
	model  string
}

func NewGeminiService(model string) (*GeminiService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return nil, err
	}

	return &GeminiService{
		client: client,
		model:  model,
	}, nil
}

func (s *GeminiService) Generate(prompt string, temperature float32, maxTokens int) (string, error) {
	ctx := context.Background()
	model := s.client.GenerativeModel(s.model)
	model.SetTemperature(temperature)
	model.SetMaxOutputTokens(maxTokens)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	var content string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					content += string(txt)
				}
			}
		}
	}

	return content, nil
}
