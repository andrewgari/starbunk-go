package tagger

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andrewgari/starbunk-go/internal/llm"
)

// Addressee represents who the message is directed to.
type Addressee string

const (
	AddresseeSelf      Addressee = "self"
	AddresseeOtherUser Addressee = "other-user"
	AddresseeRoom      Addressee = "room"
)

// Intent represents the type of message.
type Intent string

const (
	IntentQuestion  Intent = "question"
	IntentStatement Intent = "statement"
	IntentLowEffort Intent = "low-effort"
)

// StructuralTags are exact-matched enums defined by us.
type StructuralTags struct {
	Addressee Addressee `json:"addressee"`
	Intent    Intent    `json:"intent"`
}

// Result is the output of the tagger.
type Result struct {
	TopicalTags []string       `json:"topical_tags"`
	Structural  StructuralTags `json:"structural"`
}

// TaggingContext provides conversational context to the tagger.
type TaggingContext struct {
	ThreadContext       string
	ActiveConversations []string
	RecentConversations []string
	CurrentlyUsedTags   []string
}

// Service performs persona-free classification on a message.
type Service interface {
	TagMessage(ctx context.Context, messageContent string, tagCtx TaggingContext) (Result, error)
}

type llmService struct {
	llm llm.Service
}

// NewService creates a new tagger service.
func NewService(l llm.Service) Service {
	return &llmService{llm: l}
}

func (s *llmService) TagMessage(ctx context.Context, messageContent string, tagCtx TaggingContext) (Result, error) {
	systemPrompt := `You are an analytical conversation tagger. Extract topical tags and structural tags from the message.

Guidelines for topical tags:
- Topical tags should be broad concepts or specific entities (e.g., 'programming', 'binary search').
- Combine very generic and very specific tags. Example: For "I can't believe they killed off the main character in the first chapter!", tags might be: "books", "reading", "plot twist", "talking about shocking story events", "game of thrones".
- If an existing tag perfectly matches the topic, reuse it instead of creating a duplicate.
- Figure out the best tag or tags given the current conversations context.
- If the message is purely conversational ('lol', 'yeah'), topical_tags should be empty.

Addressee must be one of: 'self', 'other-user', 'room'.
Intent must be one of: 'question', 'statement', 'low-effort'.`

	if tagCtx.ThreadContext != "" {
		systemPrompt += "\n\nThread Context:\n" + tagCtx.ThreadContext
	}
	if len(tagCtx.ActiveConversations) > 0 {
		systemPrompt += "\n\nActive Conversations:"
		for _, conv := range tagCtx.ActiveConversations {
			systemPrompt += "\n- " + conv
		}
	}
	if len(tagCtx.RecentConversations) > 0 {
		systemPrompt += "\n\nRecent Conversations:"
		for _, conv := range tagCtx.RecentConversations {
			systemPrompt += "\n- " + conv
		}
	}
	if len(tagCtx.CurrentlyUsedTags) > 0 {
		systemPrompt += "\n\nCurrently Used Tags:"
		for _, tag := range tagCtx.CurrentlyUsedTags {
			systemPrompt += "\n- " + tag
		}
	}

	req := llm.GenerateRequest{
		Messages: []llm.Message{
			{Role: llm.RoleSystem, Content: systemPrompt},
			{Role: llm.RoleUser, Content: messageContent},
		},
		ExpectedOutput: llm.ResponseSchema{
			Format: llm.OutputFormatJSON,
		},
	}

	resp, err := s.llm.Generate(ctx, req)
	if err != nil {
		return Result{}, fmt.Errorf("tagger: failed to generate: %w", err)
	}

	var res Result
	if err := json.Unmarshal([]byte(resp.Text), &res); err != nil {
		return Result{}, fmt.Errorf("tagger: failed to parse json: %w", err)
	}

	return res, nil
}
