package replybot

import (
	"golang-discord-bot/internal/webhook"
	"golang-discord-bot/utils"

	"github.com/bwmarrin/discordgo"
)

type CheckBot struct {
}

func (b CheckBot) ObserverName() string {
	return "CzechBot"
}

func (b CheckBot) AvatarURL() string {
	return "https://m.media-amazon.com/images/I/21Unzn9U8sL._AC_.jpg"
}

func (b CheckBot) Response() string {
	return "I think you mean *check*"
}

func (b CheckBot) Pattern() string {
	return "\bczech\b"
}

func (b CheckBot) HandleMessage(session *discordgo.Session, message *discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) {
		webhook.WriteMessage(session, message.ChannelID, b.Response(), b.ObserverName(), b.AvatarURL())
	}
}
