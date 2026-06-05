package tagger_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/andrewgari/starbunk-go/internal/covabot/tagger"
	"github.com/andrewgari/starbunk-go/internal/llm"
)

type mockLLM struct {
	ResponseJSON string
	Error        error
	LastRequest  llm.GenerateRequest
}

func (m *mockLLM) Generate(ctx context.Context, req llm.GenerateRequest) (*llm.GenerateResponse, error) {
	m.LastRequest = req
	if m.Error != nil {
		return nil, m.Error
	}
	return &llm.GenerateResponse{
		Text: m.ResponseJSON,
	}, nil
}

func (m *mockLLM) Embed(ctx context.Context, req llm.EmbedRequest) (*llm.EmbedResponse, error) {
	return nil, nil
}

var _ = Describe("Tagger Service", func() {
	It("should parse a valid LLM response into tags", func() {
		mock := &mockLLM{
			ResponseJSON: `{"topical_tags":["programming","golang"],"structural":{"addressee":"self","intent":"question"}}`,
		}
		service := tagger.NewService(mock)

		res, err := service.TagMessage(context.Background(), "Hey Cova, what do you think of Go?", tagger.TaggingContext{})
		Expect(err).ToNot(HaveOccurred())
		Expect(res.TopicalTags).To(ConsistOf("programming", "golang"))
		Expect(res.Structural.Addressee).To(Equal(tagger.AddresseeSelf))
		Expect(res.Structural.Intent).To(Equal(tagger.IntentQuestion))
	})

	It("should handle empty topical tags", func() {
		mock := &mockLLM{
			ResponseJSON: `{"topical_tags":[],"structural":{"addressee":"room","intent":"statement"}}`,
		}
		service := tagger.NewService(mock)

		res, err := service.TagMessage(context.Background(), "Just dropping this here", tagger.TaggingContext{})
		Expect(err).ToNot(HaveOccurred())
		Expect(res.TopicalTags).To(BeEmpty())
		Expect(res.Structural.Addressee).To(Equal(tagger.AddresseeRoom))
	})
	It("should instruct the LLM to combine generic and specific tags", func() {
		mock := &mockLLM{
			ResponseJSON: `{"topical_tags":[],"structural":{"addressee":"room","intent":"statement"}}`,
		}
		service := tagger.NewService(mock)

		_, err := service.TagMessage(context.Background(), "Test message", tagger.TaggingContext{})
		Expect(err).ToNot(HaveOccurred())

		systemPrompt := mock.LastRequest.Messages[0].Content
		Expect(systemPrompt).To(ContainSubstring("Combine very generic and very specific tags"))
		Expect(systemPrompt).To(ContainSubstring("books"))
		Expect(systemPrompt).To(ContainSubstring("game of thrones"))
	})

	It("should instruct the LLM to reuse currently used tags and consider conversations", func() {
		mock := &mockLLM{
			ResponseJSON: `{"topical_tags":[],"structural":{"addressee":"room","intent":"statement"}}`,
		}
		service := tagger.NewService(mock)

		ctxData := tagger.TaggingContext{
			ActiveConversations: []string{"Technology"},
			RecentConversations: []string{"Programming"},
			CurrentlyUsedTags:   []string{"golang", "programming"},
		}

		_, err := service.TagMessage(context.Background(), "Test message", ctxData)
		Expect(err).ToNot(HaveOccurred())

		systemPrompt := mock.LastRequest.Messages[0].Content
		Expect(systemPrompt).To(ContainSubstring("Active Conversations:\n- Technology"))
		Expect(systemPrompt).To(ContainSubstring("Recent Conversations:\n- Programming"))
		Expect(systemPrompt).To(ContainSubstring("Currently Used Tags:\n- golang\n- programming"))
		Expect(systemPrompt).To(ContainSubstring("If an existing tag perfectly matches the topic, reuse it"))
	})
})
