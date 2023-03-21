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

type BotConfig struct {
	Token string
}

var starbunkConfig BotConfig
var snowbunkConfig BotConfig

func main() {
	starbunkConfig = BotConfig{Token: os.Getenv("STARBUNK_TOKEN")}
	snowbunkConfig = BotConfig{Token: os.Getenv("SNOWBUNK_TOKEN")}
	log.SetupLogger()

	starbunkClient, err := discordgo.New(starbunkConfig.Token)
	snowbunkClient, err2 := discordgo.New(snowbunkConfig.Token)

	if err != nil || err2 != nil {
		fmt.Println("Error Creating Discord Session", err)
		return
	}
	observer.MessageService = observer.MessagePublisher{Observers: make(map[string]observer.IMessageObserver)}
	observer.VoiceService = observer.VoicePublisher{Observers: make(map[string]observer.IVoiceObserver)}
	snowbunk.MessageSyncService = snowbunk.SnowbunkService{StarbunkToken: starbunkConfig.Token, SnowbunkToken: snowbunkConfig.Token}

	starbunkClient.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuildVoiceStates
	starbunkClient.AddHandler(onMessageCreate)
	starbunkClient.AddHandler(onUserVoiceStateChange)

	snowbunkClient.AddHandler(onSnowbunkMessageCreate)

	RegisterCommandBots()
	RegisterReplyBots()
	RegisterVoiceBots()

	err = starbunkClient.Open()
	if err != nil {
		fmt.Println("Error Opening Connection, ", err)
		log.ERROR.Println("Starbunk Token: ", starbunkClient.Token)
		return
	}
	err = snowbunkClient.Open()
	if err != nil {
		fmt.Println("Error Opening Snowbunk Connection, ", err)
		log.ERROR.Println("Starbunk Token: ", starbunkClient.Token)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = starbunkClient.Close()
	if err != nil {
		return
	}
	err = snowbunkClient.Close()
	if err != nil {
		return
	}
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
