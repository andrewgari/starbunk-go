package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type anthropicClient struct {
	config     ClientConfig
	httpClient *http.Client
}

type anthropicRequest struct {
	Model       string             `json:"model"`
	Messages    []anthropicMessage `json:"messages"`
	System      string             `json:"system,omitempty"`
	MaxTokens   int                `json:"max_tokens"`
	Temperature *float32           `json:"temperature,omitempty"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

func newAnthropicClient(cfg ClientConfig) Service {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.anthropic.com/v1"
	}
	return &anthropicClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *anthropicClient) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	apiReq := anthropicRequest{
		Model:       req.Model,
		MaxTokens:   4096, // A reasonable default for Anthropic
		Temperature: req.Temperature,
	}
	if apiReq.Model == "" {
		apiReq.Model = c.config.Model
	}

	var msgs []anthropicMessage
	for _, m := range req.Messages {
		if m.Role == RoleSystem {
			apiReq.System = m.Content
		} else {
			msgs = append(msgs, anthropicMessage{
				Role:    string(m.Role),
				Content: m.Content,
			})
		}
	}
	apiReq.Messages = msgs

	bodyBytes, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("anthropic: failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/messages", c.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("anthropic: failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.config.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("anthropic: request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("anthropic: unexpected status %d: %s", resp.StatusCode, string(b))
	}

	var apiResp anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("anthropic: failed to decode response: %w", err)
	}

	if len(apiResp.Content) == 0 {
		return nil, fmt.Errorf("anthropic: no content returned")
	}

	return &GenerateResponse{
		Text:             apiResp.Content[0].Text,
		PromptTokens:     apiResp.Usage.InputTokens,
		CompletionTokens: apiResp.Usage.OutputTokens,
	}, nil
}

func (c *anthropicClient) Embed(ctx context.Context, req EmbedRequest) (*EmbedResponse, error) {
	return nil, fmt.Errorf("anthropic: embed not supported")
}
