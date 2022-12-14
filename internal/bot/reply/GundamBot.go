package reply

import (
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type GundamBot struct {
	Name string
}

func (b GundamBot) ObserverName() string {
	return b.Name
}

func (b GundamBot) AvatarURL() string {
	return "https://a1.cdn.japantravel.com/photo/41317-179698/1440x960!/tokyo-unicorn-gundam-statue-in-odaiba-179698.jpg"
}

func (b GundamBot) Response() string {
	return `That's the famous Unicorn Robot, "Gandum". There, I said it.`
}

func (b GundamBot) Pattern() string {
	return "\\bg(u|a)ndam\\b"
}

func (b GundamBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) {
		webhook.WriteMessage(session, session.Identify.Token, message.ChannelID, b.Response(), b.Name, b.AvatarURL(), nil)
	}
}
