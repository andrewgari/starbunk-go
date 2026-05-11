package testenv

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// DefaultChannelID is the channel ID populated by NewMessage when WithChannelID
// is not provided. Tests that assert on CapturedMessage.ChannelID can compare
// against this constant.
const DefaultChannelID = "test-channel-1"

// DefaultGuildID is the guild ID populated by NewMessage when WithGuildID is
// not provided and the message should be a guild message. Note: NewMessage does
// NOT set a GuildID by default — use WithGuildID explicitly for guild messages.
const DefaultGuildID = "test-guild-1"

// MsgOpt is a functional option that mutates a MessageCreate.
type MsgOpt func(*discordgo.MessageCreate)

// NewMessage builds a *discordgo.MessageCreate from the given options.
// Defaults: ChannelID = DefaultChannelID, Timestamp = time.Now(), Author = &discordgo.User{}.
// GuildID is empty by default (DM context); use WithGuildID to make it a guild message.
func NewMessage(opts ...MsgOpt) *discordgo.MessageCreate {
	msg := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ID:        "test-msg-id",
			ChannelID: DefaultChannelID,
			Timestamp: time.Now(),
			Author:    &discordgo.User{ID: "test-user-id"},
		},
	}
	for _, o := range opts {
		o(msg)
	}
	return msg
}

// WithAuthorID sets the message author's Discord user ID.
func WithAuthorID(id string) MsgOpt {
	return func(m *discordgo.MessageCreate) {
		m.Author.ID = id
	}
}

// WithAuthorName sets the message author's username.
func WithAuthorName(name string) MsgOpt {
	return func(m *discordgo.MessageCreate) {
		m.Author.Username = name
	}
}

// WithBot marks the author as a bot account.
func WithBot() MsgOpt {
	return func(m *discordgo.MessageCreate) {
		m.Author.Bot = true
	}
}

// WithContent sets the message text content.
func WithContent(c string) MsgOpt {
	return func(m *discordgo.MessageCreate) {
		m.Content = c
	}
}

// WithGuildID sets the guild ID, making this a guild (non-DM) message.
func WithGuildID(id string) MsgOpt {
	return func(m *discordgo.MessageCreate) {
		m.GuildID = id
	}
}

// WithChannelID overrides the channel ID (default: DefaultChannelID).
func WithChannelID(id string) MsgOpt {
	return func(m *discordgo.MessageCreate) {
		m.ChannelID = id
	}
}

// WithAttachment adds a single placeholder attachment to the message.
func WithAttachment() MsgOpt {
	return func(m *discordgo.MessageCreate) {
		m.Attachments = append(m.Attachments, &discordgo.MessageAttachment{
			ID:       "test-attachment-id",
			Filename: "test-file.txt",
		})
	}
}

// WithTimestamp sets the message timestamp (used by OnWeekdays auditor).
func WithTimestamp(t time.Time) MsgOpt {
	return func(m *discordgo.MessageCreate) {
		m.Timestamp = t
	}
}
