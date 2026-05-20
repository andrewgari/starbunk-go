package discord

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	// webhookName is the shared name used for all Starbunk-managed webhooks.
	// All bots share one webhook per channel — this keeps usage well below
	// Discord's 15-webhook-per-channel limit.
	webhookName = "Starbunk Webhook"

	// webhookTTL is how long a channel webhook is kept alive after the last
	// send. When the TTL expires the webhook is deleted from Discord.
	webhookTTL = 5 * time.Minute

	// reaperInterval controls how often the background reaper runs.
	reaperInterval = time.Minute
)

// WebhookService executes webhook payloads on behalf of MessageService and
// manages the full lifecycle of per-channel webhooks: creation, caching, idle
// cleanup, and shutdown. This interface is intentionally unexported — callers
// interact with MessageService; WebhookService is an implementation detail.
type WebhookService interface {
	Execute(channelID, content string, id Identity) (*discordgo.Message, error)
	// Close stops the background reaper and deletes all webhooks this service
	// created. Call on bot shutdown for a clean exit.
	Close() error
}

// channelEntry holds a webhook and the last time it was used.
type channelEntry struct {
	webhook  *discordgo.Webhook
	lastUsed time.Time
}

type discordWebhookService struct {
	session *discordgo.Session
	entries map[string]*channelEntry // channelID → entry
	mu      sync.Mutex
	stopCh  chan struct{}
	doneCh  chan struct{}
}

func newDiscordWebhookService(s *discordgo.Session) WebhookService {
	ws := &discordWebhookService{
		session: s,
		entries: make(map[string]*channelEntry),
		stopCh:  make(chan struct{}),
		doneCh:  make(chan struct{}),
	}
	go ws.reaperLoop()
	return ws
}

// Execute sends content via the channel's Starbunk webhook under id.
// Identity must be valid (Username and AvatarURL non-empty); an invalid
// identity returns an error without making any Discord API calls.
// The channel's webhook entry is created lazily on the first call and its
// idle timer is reset on every successful send.
func (ws *discordWebhookService) Execute(channelID, content string, id Identity) (*discordgo.Message, error) {
	if !id.IsValid() {
		return nil, fmt.Errorf("webhook: Identity.Username and Identity.AvatarURL are required")
	}

	entry, err := ws.getOrCreate(channelID)
	if err != nil {
		return nil, fmt.Errorf("webhook: failed to get or create webhook for channel %s: %w", channelID, err)
	}

	params := &discordgo.WebhookParams{
		Content:   content,
		Username:  id.Username,
		AvatarURL: id.AvatarURL,
	}

	msg, err := ws.session.WebhookExecute(entry.webhook.ID, entry.webhook.Token, true, params)
	if err != nil {
		return nil, err
	}

	// Reset the idle timer after a confirmed successful send.
	ws.mu.Lock()
	if e, ok := ws.entries[channelID]; ok {
		e.lastUsed = time.Now()
	}
	ws.mu.Unlock()

	return msg, nil
}

// Close stops the background reaper and deletes every webhook this service
// owns. Blocks until the reaper goroutine has exited.
func (ws *discordWebhookService) Close() error {
	close(ws.stopCh)
	<-ws.doneCh
	return nil
}

// getOrCreate returns the cached entry for channelID, fetching or creating a
// webhook on Discord if none is cached. lastUsed is set to now on every hit.
func (ws *discordWebhookService) getOrCreate(channelID string) (*channelEntry, error) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if entry, ok := ws.entries[channelID]; ok {
		entry.lastUsed = time.Now()
		return entry, nil
	}

	wh, err := ws.fetchOrCreate(channelID)
	if err != nil {
		return nil, err
	}

	entry := &channelEntry{webhook: wh, lastUsed: time.Now()}
	ws.entries[channelID] = entry
	return entry, nil
}

// fetchOrCreate looks for a webhook named webhookName in channelID and returns
// it; if none exists it creates one. Called with ws.mu held.
func (ws *discordWebhookService) fetchOrCreate(channelID string) (*discordgo.Webhook, error) {
	webhooks, err := ws.session.ChannelWebhooks(channelID)
	if err != nil {
		return nil, err
	}

	for _, wh := range webhooks {
		if wh.Name == webhookName {
			return wh, nil
		}
	}

	return ws.session.WebhookCreate(channelID, webhookName, "")
}

// reaperLoop runs in its own goroutine, waking every reaperInterval to evict
// webhooks that have been idle longer than webhookTTL. On Close it deletes all
// remaining webhooks before exiting.
func (ws *discordWebhookService) reaperLoop() {
	defer close(ws.doneCh)

	ticker := time.NewTicker(reaperInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ws.reap(false)
		case <-ws.stopCh:
			ws.reap(true) // delete all on shutdown
			return
		}
	}
}

// reap deletes stale entries (idle > webhookTTL). When all is true every entry
// is deleted regardless of age — used during Close.
func (ws *discordWebhookService) reap(all bool) {
	deadline := time.Now().Add(-webhookTTL)

	ws.mu.Lock()
	defer ws.mu.Unlock()

	for channelID, entry := range ws.entries {
		if all || entry.lastUsed.Before(deadline) {
			if err := ws.session.WebhookDelete(entry.webhook.ID); err != nil {
				slog.Warn("webhook: failed to delete idle webhook",
					"channel", channelID,
					"webhook_id", entry.webhook.ID,
					"err", err,
				)
			}
			delete(ws.entries, channelID)
		}
	}
}
