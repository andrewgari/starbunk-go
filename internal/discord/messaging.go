package discord

import (
	"github.com/bwmarrin/discordgo"
)

// MessageService is the caller-facing send API. Callers declare *what* to send
// and *as whom*; the implementation decides how to deliver it (direct API vs
// webhook). Implementations are swappable — inject a stub in tests, the real
// discordMessageService in production.
type MessageService interface {
	// SendMessage sends content as the bot's own identity.
	// Use for admin, moderation, error, and ephemeral messages where
	// visuals and persona do not matter.
	SendMessage(channelID, content string) (*discordgo.Message, error)

	// SendMessageWithIdentity sends content appearing to come from id.
	// The implementation chooses the appropriate transport (webhook or
	// direct API) to honour the requested identity.
	SendMessageWithIdentity(channelID, content string, id Identity) (*discordgo.Message, error)

	// Reply sends content as the bot's own identity in reply to messageID.
	Reply(channelID, messageID, content string) (*discordgo.Message, error)

	// Edit replaces the content of an existing bot-owned message.
	Edit(channelID, messageID, content string) (*discordgo.Message, error)

	// Delete removes a message.
	Delete(channelID, messageID string) error
}

type discordMessageService struct {
	session  *discordgo.Session
	webhooks WebhookService
}

// NewMessageService returns a MessageService backed by s. SendMessageWithIdentity
// delegates to an internal WebhookService that manages per-channel webhook
// creation and caching automatically.
func NewMessageService(s *discordgo.Session) MessageService {
	return &discordMessageService{
		session:  s,
		webhooks: newDiscordWebhookService(s),
	}
}

func (ms *discordMessageService) SendMessage(channelID, content string) (*discordgo.Message, error) {
	return ms.session.ChannelMessageSend(channelID, content)
}

func (ms *discordMessageService) SendMessageWithIdentity(channelID, content string, id Identity) (*discordgo.Message, error) {
	return ms.webhooks.Execute(channelID, content, id)
}

func (ms *discordMessageService) Reply(channelID, messageID, content string) (*discordgo.Message, error) {
	return ms.session.ChannelMessageSendReply(channelID, content, &discordgo.MessageReference{
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
