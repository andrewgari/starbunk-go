package replybot

import (
	"golang-discord-bot/internal/webhook"
	"golang-discord-bot/utils"

	"github.com/bwmarrin/discordgo"
)

type SixtyNineBot struct {
}

func (b SixtyNineBot) ObserverName() string {
	return "CovaBot"
}

func (b SixtyNineBot) AvatarURL() string {
	return "https://pbs.twimg.com/profile_images/421461637325787136/0rxpHzVx.jpeg"
}

func (b SixtyNineBot) Pattern() string {
	return "\\b69|(sixty-?nine)\\b"
}

func (b SixtyNineBot) Response() string {
	return "Nice."
}

func (b SixtyNineBot) HandleMessage(session *discordgo.Session, message *discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) {
		webhook.WriteMessage(session, message.ChannelID, b.Response(), b.ObserverName(), b.AvatarURL())
	}
}
