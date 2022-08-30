package reply

import (
	"math/rand"
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"
	"time"

	"github.com/bwmarrin/discordgo"
)

type VennBot struct {
	UserID    string
	GuildID   string
	Responses []string
}

func (b VennBot) ObserverName() string {
	return ""
}

func (b VennBot) AvatarURL() string {
	return ""
}

func (b VennBot) Response() string {
	rand.Seed(time.Now().UnixNano())
	var roll = rand.Intn(len(b.Responses))
	response := b.Responses[roll]
	return response
}

func (b VennBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if message.Author.ID == b.UserID && utils.PercentChance(20) {
		var username = message.Author.Username
		var avatarURL = message.Author.AvatarURL("")
		var member, error = session.GuildMember(b.GuildID, message.Author.ID)
		if error == nil {
			memberURL := member.AvatarURL("")
			memberNick := member.Nick
			if len(memberNick) > 0 {
				username = memberNick
			}
			if len(memberURL) > 0 {
				avatarURL = memberURL
			}
		}
		webhook.WriteMessage(session, message.ChannelID, b.Response(), username, avatarURL)
	}
}
