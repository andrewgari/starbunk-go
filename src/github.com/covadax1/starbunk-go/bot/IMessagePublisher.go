package bot

import "github.com/bwmarrin/discordgo"

type IMessagePublisher interface {
	addObserver(observer IMessageObserver)
	Broadcast(session *discordgo.Session, message discordgo.Message)
}
