package replybot

import (
	"golang-discord-bot/internal/config"
	"golang-discord-bot/internal/webhook"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

type VennBot struct {
}

var id string = config.UserIDs["venn"]

const response string = "Sorry, but that was über cringe..."

func (b VennBot) ObserverName() string {
	return ""
}

func (b VennBot) AvatarURL() string {
	return ""
}

func (b VennBot) HandleMessage(session *discordgo.Session, message *discordgo.Message) {
	if message.Author.ID == id {
		rand.Seed(time.Now().UnixNano())
		var r = rand.Intn(20-1) + 1
		if r > 15 {
			var avatarUrl = message.Member.AvatarURL("")
			var username = message.Member.Nick
			webhook.WriteMessage(session, message.ChannelID, response, username, avatarUrl)
		}
	}
}
