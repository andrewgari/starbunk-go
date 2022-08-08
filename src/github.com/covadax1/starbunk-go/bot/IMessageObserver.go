package bot

import "github.com/bwmarrin/discordgo"

type IMessageObserver interface {
	ObserverName() string
	HandleMessage(session *discordgo.Session, msg discordgo.Message)
}
