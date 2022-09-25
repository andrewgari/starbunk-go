package reply

import (
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type CheckBot struct {
	Name string
}

func (b CheckBot) ObserverName() string {
	return b.Name
}

func (b CheckBot) AvatarURL() string {
	return "https://m.media-amazon.com/images/I/21Unzn9U8sL._AC_.jpg"
}

func (b CheckBot) Response() string {
	return "I think you mean *check*"
}

func (b CheckBot) Pattern() string {
	return "\\bczech\\b"
}

func (b CheckBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) {
		webhook.WriteMessage(session, message.ChannelID, b.Response(), b.Name, b.AvatarURL(), nil)
	}
}
