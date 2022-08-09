package replybot

import (
	"golang-discord-bot/internal/utils"
	"golang-discord-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type SoggyBot struct {
	Name string
	Role string
}

func (b SoggyBot) ObserverName() string {
	return b.Name
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

func (b SoggyBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) && utils.Contains(message.Member.Roles, b.Role) {
		webhook.WriteMessage(session, message.ChannelID, b.Response(), b.Name, b.AvatarURL())
	}
}
