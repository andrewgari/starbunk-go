package observer

import (
	"starbunk-bot/internal/bot/command"
	"starbunk-bot/internal/log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var MessageService IMessagePublisher
var CommandBots = make(map[string]command.ICommandBot)

type Publisher struct {
	Observers map[string]IMessageObserver
}

// addObserver implements bot.IMessagePublisher
func (p Publisher) AddObserver(observer IMessageObserver) {
	p.Observers[observer.ObserverName()] = observer
}

// broadcast implements bot.IMessagePublisher
func (p Publisher) Broadcast(session *discordgo.Session, message discordgo.Message) {
	if message.Content[:1] == "?" {
		log.INFO.Println("Got Command, ", message.Content)
		command := strings.Split(message.Content, " ")[0]
		command = strings.TrimPrefix(command, "?")
		log.INFO.Println(command)
		var bot = CommandBots[command]
		if bot != nil && bot.IsValidCommand(message.Content) {
			log.INFO.Println("Got Bot Processing")
			bot.ProcessMessage(session, message)
		}
	}
	for _, observer := range p.Observers {
		observer.HandleMessage(session, message)
	}
}
