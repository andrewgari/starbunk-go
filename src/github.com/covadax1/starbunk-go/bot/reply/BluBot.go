package reply

import (
	"log"
	"strings"

	"golang-discord-bot/src/github.com/covadax1/starbunk-go/bot/webhook"

	"github.com/bwmarrin/discordgo"
)

type BluBot struct {
	Name string
}

func (b BluBot) ObserverName() string {
	return "BluBot"
}

func (b BluBot) AvatarURL() string {
	return "https://imgur.com/WcBRCWn.png"
}

func (b BluBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	channelID := message.ChannelID
	if strings.Contains(message.Content, "blu") {
		log.Default().Println("Running BlueBot HandleMessage")
		webhook.WriteMessage(session, channelID, "Did somebody say BLU?", b.ObserverName(), b.AvatarURL())
	}
}
