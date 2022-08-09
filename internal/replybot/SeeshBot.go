package replybot

import (
	"fmt"
	"math/rand"
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type SheeshBot struct {
	Name string
	ID   string
}

func (b SheeshBot) ObserverName() string {
	return b.Name
}

func (b SheeshBot) AvatarURL() string {
	return "https://imgur.com/4CqBg7F.png"
}

func (b SheeshBot) Pattern() string {
	return "she+sh"
}

func (b SheeshBot) Response() string {
	return fmt.Sprintf("sh%ssh", strings.Repeat("e", rand.Intn(500)))
}

func (b SheeshBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if utils.Match(b.Pattern(), message.Content) || message.Author.ID == b.ID {
		webhook.WriteMessage(session, message.ChannelID, b.Response(), b.Name, b.AvatarURL())
	}
}
