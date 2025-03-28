package clients

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

var (
	ErrEmptyResponse   = errors.New("empty response from API")
	ErrAPIRequest      = errors.New("API request failed")
	ErrInvalidAPIKey   = errors.New("invalid API key")
	ErrContextCanceled = errors.New("request canceled by context")
)

type OpenAIClient struct {
	client *openai.Client
	model  string
}

func NewOpenAIClient() *OpenAIClient {
	apiKey := os.Getenv("GLHF_API_KEY")
	if apiKey == "" {
		panic("GLHF_API_KEY environment variable not set")
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.glhf.chat/v1"

	config.HTTPClient = &http.Client{
		Timeout: 30 * time.Second,
	}

	return &OpenAIClient{
		client: openai.NewClientWithConfig(config),
		model:  "hf:meta-llama/Llama-3.3-70B-Instruct",
	}
}

func (c *OpenAIClient) Chat(ctx context.Context, messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       c.model,
			Messages:    messages,
			Temperature: 0.4,
			MaxTokens:   1000,
		},
	)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", ErrEmptyResponse
	}

	return resp.Choices[0].Message.Content, nil
}
