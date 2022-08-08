package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang-discord-bot/internal/observer"
	replybot "golang-discord-bot/internal/reply-bot"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")

type Configuration struct {
	Token string
}

func readJSON() string {
	c := flag.String("c", "config.json", "Specify the configuration file.")
	flag.Parse()
	file, err := os.Open(*c)
	if err != nil {
		log.Error("can't open config file: ", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config := Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("can't decode config JSON: ", err)
	}
	return Config.Token
}

func main() {
	token := readJSON()

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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	observer.MessageService.Broadcast(s, *m.Message)
}

func registerBots() {
	var bluBot observer.IMessageObserver = replybot.BluBot{}
	observer.MessageService.AddObserver(bluBot)
}
