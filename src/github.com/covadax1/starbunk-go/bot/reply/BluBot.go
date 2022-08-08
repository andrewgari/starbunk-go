package reply

import (
	"github.com/bwmarrin/discordgo"
)

type BluBot struct {
	Name string
}

func (b BluBot) ObserverName() string {
	return "BluBot"
}

func (b BluBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	channelID := message.ChannelID
	session.ChannelMessageSend(channelID, "Did Somebody Say Blu?")
}
