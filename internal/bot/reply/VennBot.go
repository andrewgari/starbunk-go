package reply

import (
	"starbunk-bot/internal/log"
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"
	"strings"

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
	response := responses[utils.RandomRoll(len(responses))]
	return response
}

func (b VennBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if message.Author.ID == b.UserID {
		log.INFO.Println("Hello I'm Venn")
		response := ""
		if utils.PercentChance(1) {
			log.INFO.Println("OOOOOO BA NA NA")
			response = b.Response(b.Bananasponses)

		} else if utils.PercentChance(5) {
			response = strings.ToUpper(response)
			mapFunction := func(r rune) rune {
				if utils.PercentChance(50) {
					return r + 32
				}
				return r
			}
			response = strings.Map(mapFunction, response)
		} else if len(response) == 0 && utils.PercentChance(20) {
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
