package bot

import (
	"github.com/bwmarrin/discordgo"
)

var MessageService IMessagePublisher

type Publisher struct {
	Observers map[string]IMessageObserver
}

// addObserver implements bot.IMessagePublisher
func (p Publisher) AddObserver(observer IMessageObserver) {
	p.Observers[observer.ObserverName()] = observer
}

// broadcast implements bot.IMessagePublisher
func (p Publisher) Broadcast(session *discordgo.Session, message discordgo.Message) {
	for _, observer := range p.Observers {
		observer.HandleMessage(session, message)
	}
}
