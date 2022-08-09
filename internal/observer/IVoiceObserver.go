package observer

import "github.com/bwmarrin/discordgo"

type IVoiceObserver interface {
	ObserverName() string
	HandleVoiceStateChange(session *discordgo.Session, event discordgo.VoiceStateUpdate)
}
