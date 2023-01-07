package reply

import (
	"fmt"
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type PickleBot struct {
	Name string
	ID   string
}

func (b PickleBot) ObserverName() string {
	return "Klavier Gavin"
	// return b.Name
}

func (b PickleBot) AvatarURL() string {
	return "https://static.wikia.nocookie.net/aceattorney/images/f/f6/Klavier_Air_Guitar_3.gif"
	// return "https://i.imgur.com/D0czJFu.jpg"
}

func (b PickleBot) Pattern() string {
	return "klavier"
}

func (b PickleBot) Response() string {
	return "What a cool thing to say; Rock on %s"
	// return "Could you repeat that? I don't speak gremlin..."
}

func (b PickleBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) || (message.Author.ID == b.ID && utils.PercentChance(20)) {
		name := message.Member.Nick
		if len(name) > 0 {
			name = message.Author.Username
		}
		webhook.WriteMessage(session, session.Identify.Token, message.ChannelID, fmt.Sprintf(b.Response(), message.Member.Nick), b.ObserverName(), b.AvatarURL(), nil)
	}
}
