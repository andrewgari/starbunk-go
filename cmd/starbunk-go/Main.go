package main

import (
	"fmt"
	"os"
	"os/signal"
	"starbunk-bot/internal/log"
	"starbunk-bot/internal/observer"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Configuration struct {
	Token string
}

func main() {
	token := os.Getenv("TOKEN")
	log.SetupLogger()

	client, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("Error Creating Discord Session", err)
		return
	}
	observer.MessageService = observer.MessagePublisher{Observers: make(map[string]observer.IMessageObserver)}
	observer.VoiceService = observer.VoicePublisher{Observers: make(map[string]observer.IVoiceObserver)}
	client.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuildVoiceStates
	client.AddHandler(onMessageCreate)
	client.AddHandler(onUserVoiceStateChange)
	RegisterCommandBots()
	RegisterReplyBots()
	RegisterVoiceBots()

	err = client.Open()
	if err != nil {
		fmt.Println("Error Opening Connection, ", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	client.Close()
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	observer.MessageService.Broadcast(s, *m.Message)
}

func onUserVoiceStateChange(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if s.State.User.ID == v.UserID {
		return
	}
	observer.VoiceService.Broadcast(s, *v)
}
