package discord

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// WebhookService executes Discord webhook payloads. It owns the full lifecycle
// of per-channel webhooks: creation, caching, and execution.
type WebhookService interface {
	Execute(channelID string, msg WebhookMessage) (*discordgo.Message, error)
}

type discordWebhookService struct {
	session *discordgo.Session
	cache   map[string]*discordgo.Webhook // channelID → owned webhook
	mu      sync.RWMutex
}

func newDiscordWebhookService(s *discordgo.Session) WebhookService {
	return &discordWebhookService{
		session: s,
		cache:   make(map[string]*discordgo.Webhook),
	}
}

// Execute sends msg via the channel's bot-owned webhook. It validates that
// msg.Identity is non-zero before making any Discord API calls.
func (ws *discordWebhookService) Execute(channelID string, msg WebhookMessage) (*discordgo.Message, error) {
	if !msg.Identity.IsValid() {
		return nil, fmt.Errorf("webhook: Identity.Username and Identity.AvatarURL are required")
	}

	webhook, err := ws.getOrCreate(channelID)
	if err != nil {
		return nil, fmt.Errorf("webhook: failed to get or create webhook for channel %s: %w", channelID, err)
	}

	params := &discordgo.WebhookParams{
		Content:   msg.Content,
		Username:  msg.Identity.Username,
		AvatarURL: msg.Identity.AvatarURL,
	}

	return ws.session.WebhookExecute(webhook.ID, webhook.Token, true, params)
}

// getOrCreate returns the cached webhook for channelID, fetching or creating
// one if the cache is cold.
func (ws *discordWebhookService) getOrCreate(channelID string) (*discordgo.Webhook, error) {
	ws.mu.RLock()
	if wh, ok := ws.cache[channelID]; ok {
		ws.mu.RUnlock()
		return wh, nil
	}
	ws.mu.RUnlock()

	ws.mu.Lock()
	defer ws.mu.Unlock()

	// Double-check under write lock in case another goroutine just populated it.
	if wh, ok := ws.cache[channelID]; ok {
		return wh, nil
	}

	wh, err := ws.fetchOrCreate(channelID)
	if err != nil {
		return nil, err
	}

	ws.cache[channelID] = wh
	return wh, nil
}

func (ws *discordWebhookService) fetchOrCreate(channelID string) (*discordgo.Webhook, error) {
	webhooks, err := ws.session.ChannelWebhooks(channelID)
	if err != nil {
		return nil, err
	}

	for _, wh := range webhooks {
		if wh.User != nil && wh.User.ID == ws.session.State.User.ID {
			return wh, nil
		}
	}

	return ws.session.WebhookCreate(channelID, "Starbunk Webhook", "")
}
