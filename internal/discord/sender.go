package discord

import "github.com/bwmarrin/discordgo"

// MessageSender defines the strategy for sending a message to Discord.
type MessageSender interface {
	SendMessage(channelID, content string) (*discordgo.Message, error)
}

// SDKMessageSender implements MessageSender using the standard discordgo Session.
type SDKMessageSender struct {
	Session *discordgo.Session
}

// NewSDKMessageSender creates a new SDKMessageSender.
func NewSDKMessageSender(s *discordgo.Session) *SDKMessageSender {
	return &SDKMessageSender{Session: s}
}

// SendMessage sends a message using the discordgo SDK.
func (s *SDKMessageSender) SendMessage(channelID, content string) (*discordgo.Message, error) {
	return s.Session.ChannelMessageSend(channelID, content)
}
