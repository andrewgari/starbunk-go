package reply

import (
	"fmt"
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type EzioBot struct {
	Name string
	ID   string
}

func (b EzioBot) ObserverName() string {
	return b.Name
}

func (b EzioBot) AvatarURL() string {
	return "https://miro.medium.com/max/1838/1*CXPsg1BV8fuPUKchM6Cp-A.png"
}

func (b EzioBot) Response() string {
	return `Nothing is true; everything is permitted. <@%s>`
}

func (b EzioBot) Pattern() string {
	return "\\bezio|h?assassin.*\\b"
}

func (b EzioBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if message.Author.ID == b.ID && utils.Match(b.Pattern(), message.Content) {
		webhook.WriteMessage(session, session.Identify.Token, message.ChannelID, fmt.Sprintf(b.Response(), b.ID), b.Name, b.AvatarURL(), nil)
	}
}
