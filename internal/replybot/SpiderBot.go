package replybot

import (
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type SpiderBot struct {
	Name string
}

func (b SpiderBot) ObserverName() string {
	return b.Name
}

func (b SpiderBot) AvatarURL() string {
	return "https://i.dlpng.com/static/png/6569125_preview.png"
}

func (b SpiderBot) Pattern() string {
	return "\\bspider\\s?man\\b"
}

func (b SpiderBot) Response() string {
	return `Hey, it's "Spider-Man"! Don't forget the hyphen!`
}

func (b SpiderBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) {
		webhook.WriteMessage(session, message.ChannelID, b.Response(), b.Name, b.AvatarURL())
	}
}
