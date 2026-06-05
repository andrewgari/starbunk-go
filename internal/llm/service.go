package llm

import "context"

// Service is the shared abstraction for bots to interact with an LLM.
type Service interface {
	// Generate sends a prompt request to the LLM and returns the generated response.
	Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)

	// Embed generates vector embeddings for a given input text.
	Embed(ctx context.Context, req EmbedRequest) (*EmbedResponse, error)
}

// Registry provides access to different tiers of LLM services.
type Registry interface {
	High() Service
	Medium() Service
	Low() Service
}
