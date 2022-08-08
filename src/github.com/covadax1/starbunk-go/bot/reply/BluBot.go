package bot

import "github.com/bwmarrin/discordgo"

type BluBot struct {
	name string
}

func init() {
	// MessagePublisher.addObserver(BluBot{name: "BluBot"})
}

func (b BluBot) Name() string {
	return "BluBot"
}

func (b BluBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	channelID := message.ChannelID
	session.ChannelMessageSend(channelID, "Did Somebody Say Blu?")
}
