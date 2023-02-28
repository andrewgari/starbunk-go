package reply

import (
	"log"
	"math/rand"
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"
	"time"

	"github.com/bwmarrin/discordgo"
)

type QuoteBot struct {
	GuildID   string
	Responses map[string][]string
}

func (b QuoteBot) ObserverName() string {
	return "QuoteBot"
}

func (b QuoteBot) AvatarURL() string {
	return ""
}

func (b QuoteBot) Response(responses []string) string {
	rand.Seed(time.Now().UnixNano())
	var roll = rand.Intn(len(responses))
	response := responses[roll]
	return response
}

func (b QuoteBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	responses, exists := b.Responses[message.Author.ID]
	if exists && utils.PercentChance(20) {
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
		log.Printf("got meessage")
		webhook.WriteMessage(session, session.Identify.Token, message.ChannelID, b.Response(responses), username, avatarURL, nil)
	}
}
