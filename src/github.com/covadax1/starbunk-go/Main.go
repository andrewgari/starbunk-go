package main

import (
	"fmt"

	"golang-discord-bot/src/github.com/covadax1/starbunk-go/bot"

	"github.com/bwmarrin/discordgo"
)

func init() {
	CreatePublisher()
}

func CreatePublisher() {
	bot.MessageService = bot.Publisher{Observers: make(map[string]bot.IMessageObserver)}
}

func main() {
	var token = ""
	bot, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("Error Creating Discord Session", err)
		return
	}

	bot.AddHandler(messageCreate)

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	bot.MessageService.Broadcast(s, *m.Message)
}
