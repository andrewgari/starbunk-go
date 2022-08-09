package replybot

import (
	"golang-discord-bot/internal/utils"
	"golang-discord-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type GundamBot struct {
}

func (b GundamBot) ObserverName() string {
	return "That Famous Unicorn Robot, \"Gandum\""
}

func (b GundamBot) AvatarURL() string {
	return "https://a1.cdn.japantravel.com/photo/41317-179698/1440x960!/tokyo-unicorn-gundam-statue-in-odaiba-179698.jpg"
}

func (b GundamBot) Response() string {
	return `That's the famous Unicorn Robot, "Gandum". There, I said it.`
}

func (b GundamBot) Pattern() string {
	return "\bg(u|a)ndam\b"
}

func (b GundamBot) HandleMessage(session *discordgo.Session, message *discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) {
		webhook.WriteMessage(session, message.ChannelID, b.Response(), b.ObserverName(), b.AvatarURL())
	}
}
