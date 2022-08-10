package voice

import (
	"starbunk-bot/internal/log"
	"starbunk-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

var avatarUrl string = "https://i.imgur.com/v3E8yWY.jpg"

type GuildusBot struct {
	GuildusID       string
	WhaleWatchersID string
}

func (b GuildusBot) ObserverName() string {
	return "GuildusBot"
}

func (b GuildusBot) HandleVoiceStateChange(session *discordgo.Session, event discordgo.VoiceStateUpdate) {
	if event.UserID == b.GuildusID && event.ChannelID == b.WhaleWatchersID {
		webhook.WriteMessage(session, b.WhaleWatchersID, ":wave:", "GuildusBot", avatarUrl)
		_, err := session.ChannelMessageSend(b.WhaleWatchersID, ":wave:")
		if err != nil {
			log.ERROR.Println("Error Waving at Guildus.")
		}
	}
}
