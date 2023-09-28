package reply

import (
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type TheBabyBot struct {
	Name string
}

func (b TheBabyBot) ObserverName() string {
	return b.Name
}

func (b TheBabyBot) AvatarURL() string {
	return "https://i.redd.it/qc9qus78dc581.jpg"
}

func (b TheBabyBot) Response() string {
	return "https://media.tenor.com/NpnXNhWqKcwAAAAC/metroid-samus-aran.gif"
}

func (b TheBabyBot) Pattern() string {
	return "\\bbaby\\b"
}

func (b TheBabyBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) {
		webhook.WriteMessage(session, session.Identify.Token, message.ChannelID, b.Response(), b.Name, b.AvatarURL(), nil)
	}
}
