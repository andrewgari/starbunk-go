package reply

import (
	"math/rand"
	"starbunk-bot/internal/log"
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"
	"time"

	"github.com/bwmarrin/discordgo"
)

type BananaBot struct {
	UserID    string
	GuildID   string
	Responses []string
}

func (b BananaBot) ObserverName() string {
	return "VennBot"
}

func (b BananaBot) AvatarURL() string {
	return ""
}

func (b BananaBot) Response() string {
	rand.Seed(time.Now().UnixNano())
	var roll = rand.Intn(len(b.Responses))
	response := b.Responses[roll]
	return response
}

func (b BananaBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	log.INFO.Println("Hello I'm Venn")
	if message.Author.ID == b.UserID && utils.PercentChance(5) {
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
		webhook.WriteMessage(session, session.Identify.Token, message.ChannelID, b.Response(), username, avatarURL, nil)
	}
}
