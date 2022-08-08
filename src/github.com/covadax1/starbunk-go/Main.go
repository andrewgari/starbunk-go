package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang-discord-bot/src/github.com/covadax1/starbunk-go/bot"

	"golang-discord-bot/src/github.com/covadax1/starbunk-go/bot/reply"

	"github.com/bwmarrin/discordgo"
)

// var token = "NzU3NzM2ODk3MzgzNDk3OTE5.GixRil.TI5ko3YgPFszBkr0OoPsY3V60MdWDloc_YmFDQ"

type Configuration struct {
	Token string
}

func init() {
	bot.MessageService = bot.Publisher{Observers: make(map[string]bot.IMessageObserver)}
}

func readJSON() string {
	c := flag.String("c", "config.json", "Specify the configuration file.")
	flag.Parse()
	file, err := os.Open(*c)
	if err != nil {
		log.Fatal("can't open config file: ", err)
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

	bot, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("Error Creating Discord Session", err)
		return
	}

	bot.AddHandler(messageCreate)
	registerBots()
	bot.Identify.Intents = discordgo.IntentsGuildMessages

	err = bot.Open()
	if err != nil {
		fmt.Println("Error Opening Connection, ", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.Close()

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	bot.MessageService.Broadcast(s, *m.Message)
}

func registerBots() {
	var bluBot bot.IMessageObserver = reply.BluBot{Name: "BluBot"}
	bot.MessageService.AddObserver(bluBot)
}
