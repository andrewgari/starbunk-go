package replybot

import (
	"golang-discord-bot/internal/config"
	"golang-discord-bot/internal/webhook"
	"golang-discord-bot/utils"

	"github.com/bwmarrin/discordgo"
)

type SoggyBot struct {
}

func (b SoggyBot) id() string {
	return config.RoleIDs["WetBread"]
}

func (b SoggyBot) ObserverName() string {
	return "SoggyBot"
}

func (b SoggyBot) AvatarURL() string {
	return "https://imgur.com/OCB6i4x.jpg"
}

func (b SoggyBot) Pattern() string {
	return "wet bread"
}

func (b SoggyBot) Response() string {
	return "Sounds like somebody enjoys Wet Bread."
}

func (b SoggyBot) HandleMessage(session *discordgo.Session, message *discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) && utils.Contains(message.Member.Roles, b.id()) {
		webhook.WriteMessage(session, message.ChannelID, b.Response(), b.ObserverName(), b.AvatarURL())
	}
}
