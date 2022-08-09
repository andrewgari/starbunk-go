package observer

import "github.com/bwmarrin/discordgo"

type IVoicePublisher interface {
	AddObserver(observer IVoicePublisher)
	Broadcast(session *discordgo.Session, event discordgo.VoiceStateUpdate)
}
