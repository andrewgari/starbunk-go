package replybot

import (
	"fmt"
	"golang-discord-bot/internal/config"
	"golang-discord-bot/internal/utils"
	"golang-discord-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type EzioBot struct {
}

func (b EzioBot) ObserverName() string {
	return "Ezio Auditore Da Firenze"
}

func (b EzioBot) AvatarURL() string {
	return "https://miro.medium.com/max/1838/1*CXPsg1BV8fuPUKchM6Cp-A.png"
}

func (b EzioBot) Response() string {
	return `Nothing is true; everything is permitted. <@%s>`
}

func (b EzioBot) Pattern() string {
	return "\bezio|h?assassin.*\b"
}

func (b EzioBot) HandleMessage(session *discordgo.Session, message *discordgo.Message) {
	if message.Author.ID == config.UserIDs["Bender"] && utils.Match(b.Pattern(), message.Content) {
		webhook.WriteMessage(session, message.ChannelID, fmt.Sprintf(b.Response(), config.UserIDs["Bender"]), b.ObserverName(), b.AvatarURL())
	}
}
