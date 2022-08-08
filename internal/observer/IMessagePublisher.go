package observer

import "github.com/bwmarrin/discordgo"

type IMessagePublisher interface {
	AddObserver(observer IMessageObserver)
	Broadcast(session *discordgo.Session, message discordgo.Message)
}
