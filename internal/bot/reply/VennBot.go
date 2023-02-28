package reply

import (
	"math/rand"
	"starbunk-bot/internal/log"
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"
	"time"

	"github.com/bwmarrin/discordgo"
)

type VennBot struct {
	UserID        string
	GuildID       string
	Responses     []string
	Bananasponses []string
}

func (b VennBot) ObserverName() string {
	return "VennBot"
}

func (b VennBot) AvatarURL() string {
	return ""
}

func (b VennBot) Pattern() string {
	return "\\bcringe\\b"
}

func (b VennBot) Response(responses []string) string {
	rand.Seed(time.Now().UnixNano())
	var roll = rand.Intn(len(responses))
	response := responses[roll]
	return response
}

func (b VennBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if message.Author.ID == b.UserID {
		log.INFO.Println("Hello I'm Venn")
		response := ""
		if utils.PercentChance(1) {
			log.INFO.Println("OOOOOO BA NA NA")
			response = b.Response(b.Bananasponses)
		}
		if len(response) == 0 && utils.PercentChance(20) {
			log.INFO.Println("Venn is cringe")
			response = b.Response(b.Responses)
		}
		if len(response) > 0 {
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
			webhook.WriteMessage(session, session.Identify.Token, message.ChannelID, response, username, avatarURL, nil)
		}
	}
}
