package testenv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// CapturedMessage records one outbound message the bot attempted to send.
type CapturedMessage struct {
	ChannelID string
	Content   string
	// Username and AvatarURL are non-empty only for webhook/identity sends.
	Username  string
	AvatarURL string
	SentAt    time.Time
}

var (
	reChannelMessages  = regexp.MustCompile(`/channels/([^/]+)/messages$`)
	reChannelWebhooks  = regexp.MustCompile(`/channels/([^/]+)/webhooks$`)
	reWebhookExecute   = regexp.MustCompile(`/webhooks/([^/]+)/([^/?]+)`)
)

// FakeDiscordTransport is an http.RoundTripper that intercepts Discord REST
// calls and records what the bot tried to send, without any network I/O.
//
// Intercepted patterns:
//
//	POST /channels/{id}/messages       → captures regular sends
//	GET  /channels/{id}/webhooks       → returns empty list
//	POST /channels/{id}/webhooks       → returns a fake webhook object
//	POST /webhooks/{id}/{token}        → captures identity (webhook) sends
//	everything else                    → returns {} 200 OK
type FakeDiscordTransport struct {
	mu       sync.Mutex
	captured []CapturedMessage
}

// RoundTrip implements http.RoundTripper.
func (t *FakeDiscordTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	path := req.URL.Path

	switch {
	case req.Method == http.MethodPost && reChannelMessages.MatchString(path):
		return t.handleChannelMessageSend(req)

	case req.Method == http.MethodGet && reChannelWebhooks.MatchString(path):
		return jsonResponse(http.StatusOK, "[]")

	case req.Method == http.MethodPost && reChannelWebhooks.MatchString(path):
		m := reChannelWebhooks.FindStringSubmatch(path)
		channelID := ""
		if len(m) > 1 {
			channelID = m[1]
		}
		body := fmt.Sprintf(`{"id":"fake-wh-id","token":"fake-wh-token","channel_id":%q}`, channelID)
		return jsonResponse(http.StatusOK, body)

	case req.Method == http.MethodPost && reWebhookExecute.MatchString(path):
		return t.handleWebhookExecute(req)

	default:
		return jsonResponse(http.StatusOK, "{}")
	}
}

func (t *FakeDiscordTransport) handleChannelMessageSend(req *http.Request) (*http.Response, error) {
	m := reChannelMessages.FindStringSubmatch(req.URL.Path)
	channelID := ""
	if len(m) > 1 {
		channelID = m[1]
	}

	var payload struct {
		Content string `json:"content"`
	}
	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		_ = json.Unmarshal(body, &payload)
	}

	t.mu.Lock()
	t.captured = append(t.captured, CapturedMessage{
		ChannelID: channelID,
		Content:   payload.Content,
		SentAt:    time.Now(),
	})
	t.mu.Unlock()

	respBody := fmt.Sprintf(`{"id":"fake-msg-id","channel_id":%q,"content":%q}`, channelID, payload.Content)
	return jsonResponse(http.StatusOK, respBody)
}

func (t *FakeDiscordTransport) handleWebhookExecute(req *http.Request) (*http.Response, error) {
	// Extract channel_id from query params if present (wait=true responses include it).
	// For webhook executes we derive channelID from the webhook object rather
	// than the URL, so we record it as empty — tests can check Content instead.
	var payload struct {
		Content   string `json:"content"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
	}
	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		_ = json.Unmarshal(body, &payload)
	}

	t.mu.Lock()
	t.captured = append(t.captured, CapturedMessage{
		Content:   payload.Content,
		Username:  payload.Username,
		AvatarURL: payload.AvatarURL,
		SentAt:    time.Now(),
	})
	t.mu.Unlock()

	return jsonResponse(http.StatusOK, `{"id":"fake-msg-id","channel_id":""}`)
}

// Messages returns a snapshot of all captured messages in send order.
func (t *FakeDiscordTransport) Messages() []CapturedMessage {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]CapturedMessage, len(t.captured))
	copy(out, t.captured)
	return out
}

// Reset clears all captured messages. Called automatically by BotHarness.Run
// before each scenario.
func (t *FakeDiscordTransport) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.captured = t.captured[:0]
}

func jsonResponse(status int, body string) (*http.Response, error) {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    &http.Request{},
	}, nil
}

// ensure FakeDiscordTransport satisfies http.RoundTripper at compile time
var _ http.RoundTripper = (*FakeDiscordTransport)(nil)

// jsonResponseBytes is used internally when we need to return a body from a
// pre-marshalled struct rather than a raw string.
func jsonResponseBytes(status int, v any) (*http.Response, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(data)),
		Request:    &http.Request{},
	}, nil
}

// suppress unused warning — jsonResponseBytes is available for future use
var _ = jsonResponseBytes
