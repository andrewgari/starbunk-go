package replybot

import (
	"golang-discord-bot/internal/config"
	"golang-discord-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type VennBot struct {
}

func (b VennBot) id() string {
	return config.UserIDs["Venn"]
}

func (b VennBot) ObserverName() string {
	return ""
}

func (b VennBot) AvatarURL() string {
	return ""
}

func (b VennBot) Response() string {
	return "Sorry, but that was über cringe..."
}

func (b VennBot) HandleMessage(session *discordgo.Session, message *discordgo.Message) {
	if message.Author.ID == b.id() && roll20(15) {
		var avatarUrl = message.Member.AvatarURL("")
		var username = message.Member.Nick
		webhook.WriteMessage(session, message.ChannelID, b.Response(), username, avatarUrl)
	}
}
