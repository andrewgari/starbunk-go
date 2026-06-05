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

type googleClient struct {
	config     ClientConfig
	httpClient *http.Client
}

type googleRequest struct {
	SystemInstruction *googleContent          `json:"system_instruction,omitempty"`
	Contents          []googleContent         `json:"contents"`
	GenerationConfig  *googleGenerationConfig `json:"generationConfig,omitempty"`
}

type googleContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []googlePart `json:"parts"`
}

type googlePart struct {
	Text string `json:"text"`
}

type googleGenerationConfig struct {
	Temperature *float32 `json:"temperature,omitempty"`
}

type googleResponse struct {
	Candidates []struct {
		Content googleContent `json:"content"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
	} `json:"usageMetadata"`
}

func newGoogleClient(cfg ClientConfig) Service {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://generativelanguage.googleapis.com/v1beta"
	}
	return &googleClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *googleClient) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	apiReq := googleRequest{}

	if req.Temperature != nil {
		apiReq.GenerationConfig = &googleGenerationConfig{
			Temperature: req.Temperature,
		}
	}

	var contents []googleContent
	for _, m := range req.Messages {
		if m.Role == RoleSystem {
			apiReq.SystemInstruction = &googleContent{
				Parts: []googlePart{{Text: m.Content}},
			}
		} else {
			role := "user"
			if m.Role == RoleAssistant {
				role = "model"
			}
			contents = append(contents, googleContent{
				Role:  role,
				Parts: []googlePart{{Text: m.Content}},
			})
		}
	}
	apiReq.Contents = contents

	bodyBytes, err := json.Marshal(apiReq)
	if err != nil {
		return nil, fmt.Errorf("google: failed to marshal request: %w", err)
	}

	model := req.Model
	if model == "" {
		model = c.config.Model
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.config.BaseURL, model, c.config.APIKey)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("google: failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("google: request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google: unexpected status %d: %s", resp.StatusCode, string(b))
	}

	var apiResp googleResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("google: failed to decode response: %w", err)
	}

	if len(apiResp.Candidates) == 0 || len(apiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("google: no content returned")
	}

	return &GenerateResponse{
		Text:             apiResp.Candidates[0].Content.Parts[0].Text,
		PromptTokens:     apiResp.UsageMetadata.PromptTokenCount,
		CompletionTokens: apiResp.UsageMetadata.CandidatesTokenCount,
	}, nil
}

func (c *googleClient) Embed(ctx context.Context, req EmbedRequest) (*EmbedResponse, error) {
	return nil, fmt.Errorf("google: embed not implemented yet")
}
