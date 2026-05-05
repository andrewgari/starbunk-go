package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// MessagingService defines operations for interacting with Discord messages.
type MessagingService interface {
	SendMessage(channelID, content string) (*discordgo.Message, error)
	SendComplexMessage(channelID string, data *discordgo.MessageSend) (*discordgo.Message, error)
	ReplyMessage(channelID, messageID, content string) (*discordgo.Message, error)
	SendMessageWithIdentity(channelID, content, username, avatarURL string) (*discordgo.Message, error)
	EditMessage(channelID, messageID, content string) (*discordgo.Message, error)
	DeleteMessage(channelID, messageID string) error
}

type messagingService struct {
	session *discordgo.Session
}

// NewMessagingService creates a new MessagingService using the provided discordgo Session.
func NewMessagingService(s *discordgo.Session) MessagingService {
	return &messagingService{session: s}
}

func (s *messagingService) SendMessage(channelID, content string) (*discordgo.Message, error) {
	return s.session.ChannelMessageSend(channelID, content)
}

func (s *messagingService) SendComplexMessage(channelID string, data *discordgo.MessageSend) (*discordgo.Message, error) {
	return s.session.ChannelMessageSendComplex(channelID, data)
}

func (s *messagingService) ReplyMessage(channelID, messageID, content string) (*discordgo.Message, error) {
	return s.session.ChannelMessageSendReply(channelID, content, &discordgo.MessageReference{
		MessageID: messageID,
		ChannelID: channelID,
	})
}

func (s *messagingService) EditMessage(channelID, messageID, content string) (*discordgo.Message, error) {
	return s.session.ChannelMessageEdit(channelID, messageID, content)
}

func (s *messagingService) DeleteMessage(channelID, messageID string) error {
	return s.session.ChannelMessageDelete(channelID, messageID)
}

func (s *messagingService) SendMessageWithIdentity(channelID, content, username, avatarURL string) (*discordgo.Message, error) {
	webhook, err := s.getOrCreateWebhook(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create webhook: %w", err)
	}

	params := &discordgo.WebhookParams{
		Content:   content,
		Username:  username,
		AvatarURL: avatarURL,
	}

	return s.session.WebhookExecute(webhook.ID, webhook.Token, true, params)
}

func (s *messagingService) getOrCreateWebhook(channelID string) (*discordgo.Webhook, error) {
	webhooks, err := s.session.ChannelWebhooks(channelID)
	if err != nil {
		return nil, err
	}

	for _, wh := range webhooks {
		// Just take the first webhook managed by the bot
		if wh.User != nil && wh.User.ID == s.session.State.User.ID {
			return wh, nil
		}
	}

	// Create a new webhook if none exists
	return s.session.WebhookCreate(channelID, "Starbunk Webhook", "")
}
