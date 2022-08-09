package replybot

import (
	"golang-discord-bot/internal/config"
	"golang-discord-bot/internal/utils"
	"golang-discord-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type PickleBot struct {
}

func (b PickleBot) id() string {
	return config.UserIDs["Sig"]
}

func (b PickleBot) ObserverName() string {
	return "GremlinBot"
}

func (b PickleBot) AvatarURL() string {
	return "https://i.imgur.com/D0czJFu.jpg"
}

func (b PickleBot) Pattern() string {
	return "gremlin"
}

func (b PickleBot) Response() string {
	return "Could you repeat that? I don't speak gremlin..."
}

func (b PickleBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) || (message.Author.ID == b.id() && roll20(15)) {
		webhook.WriteMessage(session, message.ChannelID, b.Response(), b.ObserverName(), b.AvatarURL())
	}
}
