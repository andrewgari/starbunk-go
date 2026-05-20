package discord

import (
	"github.com/bwmarrin/discordgo"
)

// MessageService defines the high-level messaging operations available to bots.
// Call Send for plain bot messages; call SendAs for webhook persona messages.
// Implementations are swappable — inject a stub in tests, the real
// discordMessageService in production.
type MessageService interface {
	Send(channelID string, msg DirectMessage) (*discordgo.Message, error)
	SendAs(channelID string, msg WebhookMessage) (*discordgo.Message, error)
	Reply(channelID, messageID string, msg DirectMessage) (*discordgo.Message, error)
	Edit(channelID, messageID, content string) (*discordgo.Message, error)
	Delete(channelID, messageID string) error
}

type discordMessageService struct {
	session  *discordgo.Session
	webhooks WebhookService
}

// NewMessageService creates a MessageService backed by s. Webhook sends are
// handled by an internal WebhookService that manages per-channel webhook
// creation and caching automatically.
func NewMessageService(s *discordgo.Session) MessageService {
	return &discordMessageService{
		session:  s,
		webhooks: newDiscordWebhookService(s),
	}
}

func (ms *discordMessageService) Send(channelID string, msg DirectMessage) (*discordgo.Message, error) {
	return ms.session.ChannelMessageSend(channelID, msg.Content)
}

func (ms *discordMessageService) SendAs(channelID string, msg WebhookMessage) (*discordgo.Message, error) {
	return ms.webhooks.Execute(channelID, msg)
}

func (ms *discordMessageService) Reply(channelID, messageID string, msg DirectMessage) (*discordgo.Message, error) {
	return ms.session.ChannelMessageSendReply(channelID, msg.Content, &discordgo.MessageReference{
		MessageID: messageID,
		ChannelID: channelID,
	})
}

func (ms *discordMessageService) Edit(channelID, messageID, content string) (*discordgo.Message, error) {
	return ms.session.ChannelMessageEdit(channelID, messageID, content)
}

func (ms *discordMessageService) Delete(channelID, messageID string) error {
	return ms.session.ChannelMessageDelete(channelID, messageID)
}
