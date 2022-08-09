package replybot

import (
	"golang-discord-bot/internal/utils"
	"golang-discord-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type HoldBot struct {
}

func (b HoldBot) ObserverName() string {
	return "HoldBot"
}

func (b HoldBot) AvatarURL() string {
	return "https://i.imgur.com/YPFGEzM.png"
}

func (b HoldBot) Response() string {
	return "Hold."
}

func (b HoldBot) Pattern() string {
	return "^Hold\\.?$"
}

func (b HoldBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) {
		webhook.WriteMessage(session, message.ChannelID, b.Response(), b.ObserverName(), b.AvatarURL())
	}
}
