package llm

import "context"

// Service is the shared abstraction for bots to interact with an LLM.
type Service interface {
	// Generate sends a prompt request to the LLM and returns the generated response.
	Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)
}
