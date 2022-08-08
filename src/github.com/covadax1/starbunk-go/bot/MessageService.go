package bot

import "github.com/bwmarrin/discordgo"

var MessageService IMessagePublisher

type Publisher struct {
	Observers map[string]IMessageObserver
}

// addObserver implements bot.IMessagePublisher
func (p Publisher) addObserver(observer IMessageObserver) {
	p.Observers[observer.Name()] = observer
}

// broadcast implements bot.IMessagePublisher
func (p Publisher) Broadcast(session *discordgo.Session, message discordgo.Message) {
	for _, observer := range p.Observers {
		observer.HandleMessage(session, message)
	}
}
