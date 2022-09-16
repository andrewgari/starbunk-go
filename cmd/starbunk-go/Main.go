package main

import (
	"fmt"
	"os"
	"os/signal"
	"starbunk-bot/internal/log"
	"starbunk-bot/internal/observer"
	"starbunk-bot/internal/snowbunk"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Configuration struct {
	StarbunkToken string
	SnowbunkToken string
}

func main() {
	starbunk_token := os.Getenv("STARBUNK_TOKEN")
	snowbunk_token := os.Getenv("SNOWBUNK_TOKEN") // treated as a separate bot so snowfall doesn't get blasts with my crap
	log.SetupLogger()

	starbunkClient, err := discordgo.New("Bot " + starbunk_token)
	snowbunkClient, err2 := discordgo.New("Bot " + snowbunk_token)

	if err != nil || err2 != nil {
		fmt.Println("Error Creating Discord Session", err)
		return
	}
	observer.MessageService = observer.MessagePublisher{Observers: make(map[string]observer.IMessageObserver)}
	observer.VoiceService = observer.VoicePublisher{Observers: make(map[string]observer.IVoiceObserver)}
	snowbunk.MessageSyncService = snowbunk.SnowbunkService{}

	starbunkClient.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuildVoiceStates
	starbunkClient.AddHandler(onMessageCreate)
	starbunkClient.AddHandler(onUserVoiceStateChange)

	snowbunkClient.AddHandler(onSnowbunkMessageCreate)

	RegisterCommandBots()
	RegisterReplyBots()
	RegisterVoiceBots()

	err = starbunkClient.Open()
	err2 = snowbunkClient.Open()
	if err != nil {
		fmt.Println("Error Opening Connection, ", err)
		return
	}
	if err2 != nil {
		fmt.Println("Error Opening Snowbunk Connection, ", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	starbunkClient.Close()
	snowbunkClient.Close()
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	observer.MessageService.Broadcast(s, *m.Message)
}

func onSnowbunkMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	snowbunk.MessageSyncService.SyncMessage(s, *m.Message)
}

func onUserVoiceStateChange(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if s.State.User.ID == v.UserID {
		return
	}
	observer.VoiceService.Broadcast(s, *v)
}
