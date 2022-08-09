package replybot

import (
	"fmt"
	"golang-discord-bot/internal/config"
	"golang-discord-bot/internal/utils"
	"golang-discord-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

type MacaroniBot struct {
}

func (b MacaroniBot) ObserverName() string {
	return "MacaroniBot"
}

func (b MacaroniBot) AvatarURL() string {
	return "https://i.imgur.com/fgbH6Xf.jpg"
}

var (
	vennId           string = config.UserIDs["Venn"]
	macaroniId       string = config.RoleIDs["Macaroni"]
	vennPattern      string = "\bvenn\b"
	macaroniPattern  string = "\bmacaroni\b"
	vennResponse     string = `Correction: you mean Venn "Tyrone "The "Macaroni" Man" Johnson" Caelum`
	macaroniResponse string = "Are you trying to reach <@&%s}>"
)

func (b MacaroniBot) HandleMessage(session *discordgo.Session, message *discordgo.Message) {
	if utils.Match(vennPattern, message.Content) {
		webhook.WriteMessage(session, message.ChannelID, vennResponse, b.ObserverName(), b.AvatarURL())
	}
	if utils.Match(macaroniPattern, message.Content) {
		webhook.WriteMessage(session, message.ChannelID, fmt.Sprintf(macaroniResponse, macaroniId), b.ObserverName(), b.AvatarURL())
	}
}
