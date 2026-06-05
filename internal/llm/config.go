package llm

import (
	"fmt"
	"log/slog"
	"os"
)

type ProviderType string

const (
	ProviderOpenAI    ProviderType = "openai"
	ProviderAnthropic ProviderType = "anthropic"
	ProviderOllama    ProviderType = "ollama"
	ProviderGoogle    ProviderType = "google"
	ProviderNone      ProviderType = ""
)

// ClientConfig holds the config for a specific provider client.
type ClientConfig struct {
	Provider ProviderType
	BaseURL  string // E.g., http://localhost:11434 for Ollama or https://api.openai.com/v1
	APIKey   string
	Model    string // Default model for this tier/provider
}

// TierConfig holds the config for a specific tier of LLM.
type TierConfig struct {
	Provider ProviderType
	Model    string
}

// Config manages credentials for multiple providers and assigns them to tiers.
type Config struct {
	Providers map[ProviderType]ClientConfig
	High      TierConfig
	Medium    TierConfig
	Low       TierConfig
}

// ConfigFromEnv loads providers and tier assignments from environment variables.
func ConfigFromEnv() Config {
	cfg := Config{
		Providers: make(map[ProviderType]ClientConfig),
	}

	// Helper to load Provider Credentials
	loadProvider := func(p ProviderType, prefix string) {
		apiKey := os.Getenv(prefix + "API_KEY")
		baseURL := os.Getenv(prefix + "BASE_URL")
		if apiKey != "" || baseURL != "" || p == ProviderOllama {
			cfg.Providers[p] = ClientConfig{
				Provider: p,
				BaseURL:  baseURL,
				APIKey:   apiKey,
			}
		}
	}

	loadProvider(ProviderOpenAI, "OPENAI_")
	loadProvider(ProviderAnthropic, "ANTHROPIC_")
	loadProvider(ProviderOllama, "OLLAMA_")
	loadProvider(ProviderGoogle, "GOOGLE_")

	// Load Tiers
	cfg.High = TierConfig{
		Provider: ProviderType(os.Getenv("LLM_TIER_HIGH_PROVIDER")),
		Model:    os.Getenv("LLM_TIER_HIGH_MODEL"),
	}
	cfg.Medium = TierConfig{
		Provider: ProviderType(os.Getenv("LLM_TIER_MEDIUM_PROVIDER")),
		Model:    os.Getenv("LLM_TIER_MEDIUM_MODEL"),
	}
	cfg.Low = TierConfig{
		Provider: ProviderType(os.Getenv("LLM_TIER_LOW_PROVIDER")),
		Model:    os.Getenv("LLM_TIER_LOW_MODEL"),
	}

	return cfg
}

// NewClient creates a single Service instance from a ClientConfig.
func NewClient(cfg ClientConfig) (Service, error) {
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

type tieredRegistryImpl struct {
	high   Service
	medium Service
	low    Service
}

func (r *tieredRegistryImpl) High() Service {
	if r.high != nil {
		return r.high
	}
	if r.medium != nil {
		return r.medium
	}
	return r.low
}

func (r *tieredRegistryImpl) Medium() Service {
	if r.medium != nil {
		return r.medium
	}
	if r.high != nil {
		return r.high
	}
	return r.low
}

func (r *tieredRegistryImpl) Low() Service {
	if r.low != nil {
		return r.low
	}
	if r.medium != nil {
		return r.medium
	}
	return r.high
}

// NewRegistry creates a registry of LLM services based on the global Config.
func NewRegistry(cfg Config) (Registry, error) {
	createService := func(tierName string, t TierConfig) (Service, error) {
		if t.Provider == ProviderNone {
			return nil, nil
		}
		pcfg, ok := cfg.Providers[t.Provider]
		if !ok {
			slog.Warn("llm: credentials not found for provider, skipping tier", "tier", tierName, "provider", t.Provider)
			return nil, nil
		}

		// Override the model with the tier's model
		pcfg.Model = t.Model
		return NewClient(pcfg)
	}

	high, err := createService("high", cfg.High)
	if err != nil {
		return nil, fmt.Errorf("failed to create high tier: %w", err)
	}
	medium, err := createService("medium", cfg.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to create medium tier: %w", err)
	}
	low, err := createService("low", cfg.Low)
	if err != nil {
		return nil, fmt.Errorf("failed to create low tier: %w", err)
	}

	if high == nil && medium == nil && low == nil {
		slog.Warn("llm: no valid LLM tiers were configured. Bot may fail to generate responses.")
	}

	return &tieredRegistryImpl{
		high:   high,
		medium: medium,
		low:    low,
	}, nil
}
