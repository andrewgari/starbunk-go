package voice

import (
	"starbunk-bot/internal/log"

	"github.com/bwmarrin/discordgo"
)

type GuyChannelBot struct {
	GuysID           string
	GuysChannelID    string
	NotGuysChannelID string
	LoungeId         string
	GuildID          string
}

func (b GuyChannelBot) ObserverName() string {
	return "GuyVoiceBot"
}

func (b GuyChannelBot) HandleVoiceStateChange(session *discordgo.Session, event discordgo.VoiceStateUpdate) {
	if event.ChannelID == b.GuysChannelID && event.UserID != b.GuysID {
		b.moveUser(session, event.UserID)
	}
	if event.ChannelID == b.NotGuysChannelID && event.UserID == b.GuysID {
		b.moveUser(session, b.GuysID)
	}
}

func (b GuyChannelBot) moveUser(session *discordgo.Session, userID string) {
	err := session.GuildMemberMove(b.GuildID, userID, &b.LoungeId)
	if err != nil {
		log.ERROR.Println("Error Moving User to Another Channel")
	}
}
