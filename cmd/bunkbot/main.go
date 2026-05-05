package main

import (
	"log"

	"github.com/andrewgari/starbunk-go/internal/bot"
	"github.com/andrewgari/starbunk-go/internal/discord"
	"github.com/bwmarrin/discordgo"
)

func main() {
	bot.Run("BunkBot", messageCreate)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "ping bunkbot" {
		sender := discord.NewMessagingService(s)
		_, err := sender.SendMessage(m.ChannelID, "Pong from bunkbot!")
		if err != nil {
			log.Printf("failed to send message: %v\n", err)
		}
	}
}
