package conversation

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/andrewgari/starbunk-go/internal/llm"
)

// ActiveConversation represents an ongoing thread in working memory.
type ActiveConversation struct {
	ID           string
	Centroid     []float32
	LastActivity time.Time
	MessageCount int
}

// Tracker assigns incoming messages to conversations based on topical embeddings.
type Tracker interface {
	Assign(ctx context.Context, channelID string, tags []string) ([]string, error)
}

type llmTracker struct {
	llm   llm.Service
	mu    sync.RWMutex
	live  map[string][]*ActiveConversation
	tHigh float32
	tLow  float32
}

// NewTracker creates a new conversation tracker with default similarity thresholds.
func NewTracker(l llm.Service) Tracker {
	return &llmTracker{
		llm:   l,
		live:  make(map[string][]*ActiveConversation),
		tHigh: 0.75,
		tLow:  0.45,
	}
}

// simResult pairs a conversation with its pre-computed similarity to the incoming message.
type simResult struct {
	conv *ActiveConversation
	sim  float32
}

func (t *llmTracker) Assign(ctx context.Context, channelID string, tags []string) ([]string, error) {
	if len(tags) == 0 {
		return nil, nil
	}

	embReq := llm.EmbedRequest{Input: tags}
	embResp, err := t.llm.Embed(ctx, embReq)
	if err != nil {
		return nil, fmt.Errorf("conversation: embed failed: %w", err)
	}

	msgCentroid := averageEmbeddings(embResp.Embeddings)

	t.mu.Lock()
	defer t.mu.Unlock()

	live := t.live[channelID]
	var assigned []string
	var ambiguous []simResult // carries pre-computed sim so we don't recompute below

	for _, conv := range live {
		sim := cosineSimilarity(msgCentroid, conv.Centroid)
		if sim >= t.tHigh {
			assigned = append(assigned, conv.ID)
			updateCentroid(conv, msgCentroid)
			conv.LastActivity = time.Now()
		} else if sim >= t.tLow {
			ambiguous = append(ambiguous, simResult{conv, sim})
		}
	}

	if len(assigned) == 0 && len(ambiguous) == 0 {
		// Seed new conversation — no existing thread is close enough.
		id := "conv-" + generateID()
		newConv := &ActiveConversation{
			ID:           id,
			Centroid:     msgCentroid,
			LastActivity: time.Now(),
			MessageCount: 1,
		}
		t.live[channelID] = append(t.live[channelID], newConv)
		assigned = append(assigned, id)
	} else if len(assigned) == 0 && len(ambiguous) > 0 {
		// Ambiguous band: assign to the best match using the already-computed similarities.
		var best simResult
		best.sim = -1
		for _, r := range ambiguous {
			if r.sim > best.sim {
				best = r
			}
		}
		if best.conv != nil {
			assigned = append(assigned, best.conv.ID)
			updateCentroid(best.conv, msgCentroid)
			best.conv.LastActivity = time.Now()
		}
	}

	return assigned, nil
}

func updateCentroid(c *ActiveConversation, msgCentroid []float32) {
	c.MessageCount++
	for i := range c.Centroid {
		c.Centroid[i] = c.Centroid[i] + (msgCentroid[i]-c.Centroid[i])/float32(c.MessageCount)
	}
}

func averageEmbeddings(embs [][]float32) []float32 {
	if len(embs) == 0 {
		return nil
	}
	dim := len(embs[0])
	res := make([]float32, dim)
	valid := 0
	for _, emb := range embs {
		if len(emb) != dim {
			// Skip malformed embeddings rather than panic on out-of-bounds access.
			continue
		}
		for i := range dim {
			res[i] += emb[i]
		}
		valid++
	}
	if valid == 0 {
		return nil
	}
	for i := range dim {
		res[i] /= float32(valid)
	}
	return res
}

func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float32
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / float32(math.Sqrt(float64(normA))*math.Sqrt(float64(normB)))
}

func generateID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
