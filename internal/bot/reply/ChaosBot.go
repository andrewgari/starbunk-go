package reply

import (
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type ChaosBot struct {
	Name string
}

func (b ChaosBot) ObserverName() string {
	return b.Name
}

func (b ChaosBot) AvatarURL() string {
	return "https://preview.redd.it/md0lzbvuc3571.png?width=1920&format=png&auto=webp&s=ff403a8d4b514af8d99792a275d2c066b8d1a4de"
}

func (b ChaosBot) Response() string {
	return "All I know is...I'm here to kill chaos"
}

func (b ChaosBot) Pattern() string {
	return "\\bchaos\\b"
}

func (b ChaosBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) {
		webhook.WriteMessage(session, session.Identify.Token, message.ChannelID, b.Response(), b.Name, b.AvatarURL(), nil)
	}
}
