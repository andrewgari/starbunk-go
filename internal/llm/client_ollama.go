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

type ollamaClient struct {
	config     ClientConfig
	httpClient *http.Client
}

type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Options  *ollamaOptions  `json:"options,omitempty"`
}

type ollamaOptions struct {
	Temperature *float32 `json:"temperature,omitempty"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	PromptEvalCount int `json:"prompt_eval_count"`
	EvalCount       int `json:"eval_count"`
}

type ollamaEmbedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type ollamaEmbedResponse struct {
	Model      string      `json:"model"`
	Embeddings [][]float32 `json:"embeddings"`
}

func newOllamaClient(cfg ClientConfig) Service {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:11434"
	}
	return &ollamaClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *ollamaClient) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	apiReq := ollamaRequest{
		Model:    req.Model,
		Stream:   false,
		Messages: make([]ollamaMessage, len(req.Messages)),
	}
	if apiReq.Model == "" {
		apiReq.Model = c.config.Model
	}
	if req.Temperature != nil {
		apiReq.Options = &ollamaOptions{
			Temperature: req.Temperature,
		}
	}
	for i, m := range req.Messages {
		apiReq.Messages[i] = ollamaMessage{
			Role:    string(m.Role),
			Content: m.Content,
		}
	}

	bodyBytes, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("ollama: failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/chat", c.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("ollama: failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama: request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama: unexpected status %d: %s", resp.StatusCode, string(b))
	}

	var apiResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("ollama: failed to decode response: %w", err)
	}

	return &GenerateResponse{
		Text:             apiResp.Message.Content,
		PromptTokens:     apiResp.PromptEvalCount,
		CompletionTokens: apiResp.EvalCount,
	}, nil
}

func (c *ollamaClient) Embed(ctx context.Context, req EmbedRequest) (*EmbedResponse, error) {
	model := req.Model
	if model == "" {
		model = c.config.Model
	}

	apiReq := ollamaEmbedRequest{
		Model: model,
		Input: req.Input,
	}

	bodyBytes, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("ollama: failed to marshal embed request: %w", err)
	}

	url := fmt.Sprintf("%s/api/embed", c.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("ollama: failed to create embed request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama: embed request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama: unexpected status %d: %s", resp.StatusCode, string(b))
	}

	var apiResp ollamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("ollama: failed to decode embed response: %w", err)
	}

	return &EmbedResponse{
		Embeddings: apiResp.Embeddings,
	}, nil
}
