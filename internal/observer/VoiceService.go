package observer

import "github.com/bwmarrin/discordgo"

var VoiceService IVoicePublisher

type VoicePublisher struct {
	Observers map[string]IVoiceObserver
}

func (p VoicePublisher) AddObserver(observer IVoiceObserver) {
	p.Observers[observer.ObserverName()] = observer
}

func (p VoicePublisher) Broadcast(session *discordgo.Session, event discordgo.VoiceStateUpdate) {
	for _, observer := range p.Observers {
		observer.HandleVoiceStateChange(session, event)
	}
}
