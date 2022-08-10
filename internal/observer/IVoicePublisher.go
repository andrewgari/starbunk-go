package observer

import "github.com/bwmarrin/discordgo"

type IVoicePublisher interface {
	AddObserver(observer IVoiceObserver)
	Broadcast(session *discordgo.Session, event discordgo.VoiceStateUpdate)
}
