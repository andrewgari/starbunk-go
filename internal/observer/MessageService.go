package observer

import (
	"regexp"
	"starbunk-bot/internal/bot/command"
	"starbunk-bot/internal/log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var MessageService IMessagePublisher
var CommandBots = make(map[string]command.ICommandBot)

type MessagePublisher struct {
	Observers map[string]IMessageObserver
}

// addObserver implements bot.IMessagePublisher
func (p MessagePublisher) AddObserver(observer IMessageObserver) {
	p.Observers[observer.ObserverName()] = observer
}

// broadcast implements bot.IMessagePublisher
func (p MessagePublisher) Broadcast(session *discordgo.Session, message discordgo.Message) {
	if len(message.Content) > 0 && message.Content[0:1] == "!" {
		whitespaces := regexp.MustCompile(`\s+`)
		message.Content = whitespaces.ReplaceAllString(message.Content, " ")
		log.INFO.Println("Got Command, ", message.Content)
		command := strings.Split(message.Content, " ")[0]
		command = strings.TrimPrefix(command, "!")
		log.INFO.Println(command)
		var bot = CommandBots[command]
		if bot != nil && bot.IsValidCommand(message.Content) {
			log.INFO.Println("Got Bot Processing")
			bot.ProcessMessage(session, message)
		}
	}
	for _, observer := range p.Observers {
		go observer.HandleMessage(session, message)
	}
}
