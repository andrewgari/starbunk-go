package reply

import (
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type VennBot struct {
	ID string
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
	if message.Author.ID == b.ID && utils.Roll20(15) {
		var avatarUrl = ""
		var username = ""
		var vennAsMember = message.Member
		if vennAsMember == nil {
			var vennAsUser = message.Author
			avatarUrl = vennAsUser.AvatarURL("")
			username = vennAsUser.Username
		} else {
			avatarUrl = vennAsMember.AvatarURL("")
			username = vennAsMember.Nick
		}
		webhook.WriteMessage(session, message.ChannelID, b.Response(), username, avatarUrl)
	}
}
