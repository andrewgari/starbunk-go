package replybot

import (
	"fmt"
	"golang-discord-bot/internal/config"
	"golang-discord-bot/internal/utils"
	"golang-discord-bot/internal/webhook"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type SheeshBot struct {
}

func (b SheeshBot) id() string {
	return config.UserIDs["Guy"]
}

func (b SheeshBot) ObserverName() string {
	return "SheeshBot"
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
	if utils.Match(b.Pattern(), message.Content) || message.Author.ID == b.id() {
		webhook.WriteMessage(session, message.ChannelID, b.Response(), b.ObserverName(), b.AvatarURL())
	}
}
