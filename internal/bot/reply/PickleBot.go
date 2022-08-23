package reply

import (
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type PickleBot struct {
	Name string
	ID   string
}

func (b PickleBot) ObserverName() string {
	return b.Name
}

func (b PickleBot) AvatarURL() string {
	return "https://i.imgur.com/D0czJFu.jpg"
}

func (b PickleBot) Pattern() string {
	return "gremlin"
}

func (b PickleBot) Response() string {
	return "Could you repeat that? I don't speak gremlin..."
}

func (b PickleBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) || (message.Author.ID == b.ID && utils.PercentChance(20)) {
		webhook.WriteMessage(session, message.ChannelID, b.Response(), b.Name, b.AvatarURL())
	}
}
