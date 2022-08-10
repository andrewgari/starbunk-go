package voice

import (
	"starbunk-bot/internal/log"

	"github.com/bwmarrin/discordgo"
)

type FeliBot struct {
	Name            string
	FeliID          string
	GuildID         string
	AFK_ID          string
	WhaleWatchersID string
}

func (b FeliBot) ObserverName() string {
	return "FeliBot"
}

func (b FeliBot) HandleVoiceStateChange(session *discordgo.Session, event discordgo.VoiceStateUpdate) {
	if event.UserID == b.FeliID && event.ChannelID == b.AFK_ID {
		// Good Night, Sweet Prince
		log.INFO.Println("Good Night, Sweet Prince")
		// Post Message to Whale Watchers
		b.sendMessage(session, b.WhaleWatchersID)

		// DM Feli
		feliDM, err := session.UserChannelCreate(b.FeliID)
		if err != nil {
			log.ERROR.Println("Error Finding Feli. He must have fallen asleep under a rock.")
		}
		b.sendMessage(session, feliDM.ID)

	}
}

func (b FeliBot) sendMessage(session *discordgo.Session, channelID string) {
	_, err := session.ChannelMessageSend(channelID, "Good Night, Sweet Prince")
	if err != nil {
		log.ERROR.Println("Error Tucking Feli in.")
	}
}
