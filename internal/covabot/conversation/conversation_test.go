package conversation_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/andrewgari/starbunk-go/internal/covabot/conversation"
	"github.com/andrewgari/starbunk-go/internal/llm"
)

type mockLLM struct {
	embeddings [][]float32
}

func (m *mockLLM) Generate(ctx context.Context, req llm.GenerateRequest) (*llm.GenerateResponse, error) {
	return nil, nil
}

func (m *mockLLM) Embed(ctx context.Context, req llm.EmbedRequest) (*llm.EmbedResponse, error) {
	return &llm.EmbedResponse{Embeddings: m.embeddings}, nil
}

var _ = Describe("Conversation Tracker", func() {
	var tracker conversation.Tracker
	var mock *mockLLM
	var chanID = "chan1"

	BeforeEach(func() {
		mock = &mockLLM{}
		tracker = conversation.NewTracker(mock)
	})

	It("seeds a new conversation when none exist", func() {
		mock.embeddings = [][]float32{{1.0, 0.0}}
		convIDs, err := tracker.Assign(context.Background(), chanID, []string{"kh"})
		Expect(err).ToNot(HaveOccurred())
		Expect(convIDs).To(HaveLen(1))
	})

	It("joins an existing conversation if similarity is high", func() {
		// First message
		mock.embeddings = [][]float32{{1.0, 0.0}}
		conv1, _ := tracker.Assign(context.Background(), chanID, []string{"kh"})

		// Second message with identical embedding
		mock.embeddings = [][]float32{{1.0, 0.0}}
		conv2, err := tracker.Assign(context.Background(), chanID, []string{"kh 2"})
		Expect(err).ToNot(HaveOccurred())
		Expect(conv2).To(Equal(conv1))
	})

	It("seeds a new conversation if similarity is below t_low", func() {
		// First message
		mock.embeddings = [][]float32{{1.0, 0.0}}
		conv1, _ := tracker.Assign(context.Background(), chanID, []string{"kh"})

		// Second message with orthogonal embedding
		mock.embeddings = [][]float32{{0.0, 1.0}}
		conv2, err := tracker.Assign(context.Background(), chanID, []string{"dbz"})
		Expect(err).ToNot(HaveOccurred())
		Expect(conv2).ToNot(Equal(conv1))
	})
})
