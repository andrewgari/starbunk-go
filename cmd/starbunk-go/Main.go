package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang-discord-bot/internal/log"
	"golang-discord-bot/internal/observer"
	"golang-discord-bot/internal/replybot"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Configuration struct {
	Token string
}

func main() {
	token := readJSON()
	log.SetupLogger()

	client, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("Error Creating Discord Session", err)
		return
	}
	observer.MessageService = observer.Publisher{Observers: make(map[string]observer.IMessageObserver)}
	client.AddHandler(messageCreate)
	registerBots()
	client.Identify.Intents = discordgo.IntentsGuildMessages

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

func readJSON() string {
	c := flag.String("c", "config.json", "Specify the configuration file.")
	flag.Parse()
	file, err := os.Open(*c)
	if err != nil {
		log.ERROR.Println("can't open config file: ", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config := Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.ERROR.Println("can't decode config JSON: ", err)
	}
	return Config.Token
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	observer.MessageService.Broadcast(s, *m.Message)
}

func registerBots() {
	var bluBot observer.IMessageObserver = replybot.BluBot{}
	observer.MessageService.AddObserver(bluBot)
	var chaosBot observer.IMessageObserver = replybot.ChaosBot{}
	observer.MessageService.AddObserver(chaosBot)
	var checkBot observer.IMessageObserver = replybot.CheckBot{}
	observer.MessageService.AddObserver(checkBot)
	var deafBot observer.IMessageObserver = replybot.DeafBot{}
	observer.MessageService.AddObserver(deafBot)
	var ezioBot observer.IMessageObserver = replybot.EzioBot{}
	observer.MessageService.AddObserver(ezioBot)
	var gundamBot observer.IMessageObserver = replybot.GundamBot{}
	observer.MessageService.AddObserver(gundamBot)
	var holdBot observer.IMessageObserver = replybot.HoldBot{}
	observer.MessageService.AddObserver(holdBot)
	var macaroniBot observer.IMessageObserver = replybot.MacaroniBot{}
	observer.MessageService.AddObserver(macaroniBot)
	var pickleBot observer.IMessageObserver = replybot.PickleBot{}
	observer.MessageService.AddObserver(pickleBot)
	var sheeshBot observer.IMessageObserver = replybot.SheeshBot{}
	observer.MessageService.AddObserver(sheeshBot)
	var sixtyNineBot observer.IMessageObserver = replybot.SixtyNineBot{}
	observer.MessageService.AddObserver(sixtyNineBot)
	var soggyBot observer.IMessageObserver = replybot.SoggyBot{}
	observer.MessageService.AddObserver(soggyBot)
	var spiderBot observer.IMessageObserver = replybot.SpiderBot{}
	observer.MessageService.AddObserver(spiderBot)
	var vennBot observer.IMessageObserver = replybot.VennBot{}
	observer.MessageService.AddObserver(vennBot)
}
