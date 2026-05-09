package main

import (
	"log"

	"github.com/andrewgari/starbunk-go/internal/bot"
	"github.com/andrewgari/starbunk-go/internal/discord"
	"github.com/andrewgari/starbunk-go/internal/middleware"
	"github.com/bwmarrin/discordgo"
)

var auditor = middleware.AllOf(
	middleware.NotSelf,
	middleware.NotBot,
	middleware.GuildOnly,
	middleware.HasContent,
)

func main() {
	bot.Run("BlueBot", auditor, messageCreate)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "ping bluebot" {
		sender := discord.NewMessagingService(s)
		_, err := sender.SendMessage(m.ChannelID, "Pong from bluebot!")
		if err != nil {
			log.Printf("failed to send message: %v\n", err)
		}
	}
}
