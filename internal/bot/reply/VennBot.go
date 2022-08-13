package reply

import (
	"starbunk-bot/internal/log"
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type VennBot struct {
	UserID  string
	GuildID string
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

func (b VennBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if message.Author.ID == b.UserID && utils.Roll20(15) {
		var username, avatarURL string
		var member, error = session.GuildMember(b.GuildID, b.UserID)
		if error != nil {
			log.ERROR.Println("Error getting avatar url from guild", error)
			avatarURL = message.Author.AvatarURL("")
			username = message.Author.Username
		} else {
			avatarURL = member.AvatarURL("")
			username = member.Nick
		}
		webhook.WriteMessage(session, message.ChannelID, b.Response(), username, avatarURL)
	}
}
