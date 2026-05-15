
package services

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"velora/config"
)

// OpenAIService handles interactions with the OpenAI API.
type OpenAIService struct {
	client *openai.Client
}

// NewOpenAIService creates a new OpenAIService.
func NewOpenAIService(cfg config.Config) (*OpenAIService, error) {
	if cfg.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := openai.NewClient(cfg.OpenAIAPIKey)
	return &OpenAIService{client: client}, nil
}

// Generate generates content from a text prompt.
func (s *OpenAIService) Generate(prompt, modelName string, temperature float32, maxOutputTokens int) (string, error) {
	if modelName == "" {
		modelName = openai.GPT4oMini
	}
	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: modelName,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature:      temperature,
			MaxTokens:        maxOutputTokens,
			TopP:             1,
			FrequencyPenalty: 0,
			PresencePenalty:  0,
		},
	)
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateStream generates content from a text prompt and streams the response.
func (s *OpenAIService) GenerateStream(prompt, modelName string, temperature float32, maxOutputTokens int) (<-chan string, <-chan error) {
	out := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errChan)

		stream, err := s.client.CreateChatCompletionStream(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: modelName,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
				Temperature:      temperature,
				MaxTokens:        maxOutputTokens,
				TopP:             1,
				FrequencyPenalty: 0,
				PresencePenalty:  0,
				Stream:           true,
			},
		)
		if err != nil {
			errChan <- fmt.Errorf("OpenAI stream error: %w", err)
			return
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if err != nil {
				errChan <- fmt.Errorf("stream receive error: %w", err)
				return
			}

			if len(response.Choices) > 0 && response.Choices[0].FinishReason == openai.FinishReasonStop {
				break
			}

			if len(response.Choices) > 0 {
				out <- response.Choices[0].Delta.Content
			}
		}
	}()

	return out, errChan
}

// Close is a no-op for the OpenAI service.
func (s *OpenAIService) Close() error {
	return nil
}
