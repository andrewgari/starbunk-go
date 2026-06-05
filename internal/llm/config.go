package llm

import (
	"fmt"
	"os"
)

type ProviderType string

const (
	ProviderOpenAI    ProviderType = "openai"
	ProviderAnthropic ProviderType = "anthropic"
	ProviderOllama    ProviderType = "ollama"
	ProviderGoogle    ProviderType = "google"
)

type Config struct {
	Provider ProviderType
	BaseURL  string // E.g., http://localhost:11434 for Ollama or https://api.openai.com/v1
	APIKey   string
	Model    string // Default model for this provider
}

// ConfigFromEnv loads the LLM configuration from environment variables.
// It uses a prefix (e.g., "CLOUD_LLM_" or "LOCAL_LLM_") to fetch the settings.
func ConfigFromEnv(prefix string) Config {
	return Config{
		Provider: ProviderType(os.Getenv(prefix + "PROVIDER")),
		BaseURL:  os.Getenv(prefix + "BASE_URL"),
		APIKey:   os.Getenv(prefix + "API_KEY"),
		Model:    os.Getenv(prefix + "MODEL"),
	}
}

// NewService creates a Service based on the provided config.
func NewService(cfg Config) (Service, error) {
	switch cfg.Provider {
	case ProviderOpenAI:
		return newOpenAIClient(cfg), nil
	case ProviderOllama:
		return newOllamaClient(cfg), nil
	case ProviderAnthropic:
		return newAnthropicClient(cfg), nil
	case ProviderGoogle:
		return newGoogleClient(cfg), nil
	default:
		return nil, fmt.Errorf("llm: unknown or unsupported provider '%s'", cfg.Provider)
	}
}
