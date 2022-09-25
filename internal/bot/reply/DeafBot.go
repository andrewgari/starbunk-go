package reply

import (
	"starbunk-bot/internal/webhook"
	"time"

	"github.com/bwmarrin/discordgo"
)

type DeafBot struct {
	Name string
	ID   string
}

func (b DeafBot) ObserverName() string {
	return b.Name
}

func (b DeafBot) AvatarURL() string {
	return "https://www.reptilecentre.com/blog/wp-content/uploads/2020/02/leopard-gecko-header.jpg"
}

func (b DeafBot) Response() string {
	return `He is **Awake** https://giphy.com/gifs/come-at-me-im-here-big-bird-H6W9H29kVsUI2hJE90`
}

var lastResponse = time.Unix(0, 0)

func (b DeafBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if message.Author.ID == b.ID && time.Now().After(lastResponse.AddDate(0, 0, 10)) {
		if !lastResponse.IsZero() {
			webhook.WriteMessage(session, message.ChannelID, b.Response(), b.Name, b.AvatarURL(), nil)
		}
		lastResponse = time.Now()
	}
}
