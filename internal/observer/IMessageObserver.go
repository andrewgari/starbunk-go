package observer

import "github.com/bwmarrin/discordgo"

type IMessageObserver interface {
	ObserverName() string
	AvatarURL() string
	HandleMessage(session *discordgo.Session, msg discordgo.Message)
}
