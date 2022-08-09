package main

import (
	"encoding/json"
	"flag"
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
	token := readJSON()
	log.SetupLogger()

	client, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("Error Creating Discord Session", err)
		return
	}
	observer.MessageService = observer.Publisher{Observers: make(map[string]observer.IMessageObserver)}
	client.AddHandler(messageCreate)
	RegisterCommandBots()
	RegisterReplyBots()
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
