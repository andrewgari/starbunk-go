package reply

import (
	"fmt"
	"starbunk-bot/internal/utils"
	"starbunk-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type MacaroniBot struct {
	Name string
	ID   string
	Role string
}

func (b MacaroniBot) ObserverName() string {
	return b.Name
}

func (b MacaroniBot) AvatarURL() string {
	return "https://i.imgur.com/fgbH6Xf.jpg"
}

var (
	vennPattern      string = "\\bvenn\\b"
	macaroniPattern  string = "\\bmacaroni\\b"
	vennResponse     string = `Correction: you mean Venn "Tyrone "The "Macaroni" Man" Johnson" Caelum`
	macaroniResponse string = "Are you trying to reach <@&%s>"
)

func (b MacaroniBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if utils.Match(vennPattern, message.Content) {
		webhook.WriteMessage(session, message.ChannelID, vennResponse, b.Name, b.AvatarURL())
	}
	if utils.Match(macaroniPattern, message.Content) {
		webhook.WriteMessage(session, message.ChannelID, fmt.Sprintf(macaroniResponse, b.ID), b.Name, b.AvatarURL())
	}
}
