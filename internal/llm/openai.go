package llm

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

// OpenAIClient implements LLMService using the official OpenAI Go SDK.
type OpenAIClient struct {
	client openai.Client
}

// NewOpenAIClient creates a new OpenAI client with proper configuration.
// It uses the official OpenAI Go SDK which handles retries, timeouts, and error handling automatically.
// The client is configured with:
// - Custom HTTP client with 30s timeout for LLM requests
// - Automatic retries (2 by default) for transient errors
// - Typed model constant ChatModelGPT4oMini
// The apiKey parameter is required and should be provided from config.Load().
func NewOpenAIClient(apiKey string) *OpenAIClient {

	// Create HTTP client with appropriate timeout for LLM requests
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create client with API key and custom HTTP client
	// The SDK will use this HTTP client for all requests
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithHTTPClient(httpClient),
		// Configure retries: 2 retries by default for transient errors
		// (408 Request Timeout, 409 Conflict, 429 Rate Limit, >=500 Internal errors)
		option.WithMaxRetries(2),
	)

	return &OpenAIClient{
		client: client,
	}
}

// Ask sends a prompt to the OpenAI API and returns the response.
// It uses the official OpenAI SDK which handles HTTP requests, retries, and error handling.
// The request includes:
// - System message with instructions
// - User prompt
// - Typed model constant ChatModelGPT4oMini
// - Automatic retries on transient errors
func (c *OpenAIClient) Ask(ctx context.Context, prompt string) (string, error) {
	if prompt == "" {
		return "", fmt.Errorf("prompt cannot be empty")
	}

	// Create chat completion request using the official SDK
	// Using typed model constant ChatModelGPT4oMini and helper functions for messages
	chatCompletion, err := c.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model: shared.ChatModelGPT4oMini,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("Eres un agente de clima experto en Nicaragua. Responde solo JSON."),
				openai.UserMessage(prompt),
			},
		},
		// Optionally configure per-request retries (overrides client default)
		// option.WithMaxRetries(3),
	)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	// Extract the response content
	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("openai API returned empty response")
	}

	// Get the content from the first choice
	content := chatCompletion.Choices[0].Message.Content
	if content == "" {
		return "", fmt.Errorf("openai API returned empty content")
	}

	return content, nil
}
