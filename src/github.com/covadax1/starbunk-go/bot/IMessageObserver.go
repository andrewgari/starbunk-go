package bot

import "github.com/bwmarrin/discordgo"

type IMessageObserver interface {
	Name() string
	HandleMessage(session *discordgo.Session, msg discordgo.Message)
}
