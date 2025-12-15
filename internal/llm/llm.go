package llm

import "context"

// LLMService defines the interface for language model interactions.
type LLMService interface {
	Ask(ctx context.Context, prompt string) (string, error)
}
