package main

import (
	"log/slog"

	"github.com/andrewgari/starbunk-go/internal/bot"
	"github.com/andrewgari/starbunk-go/internal/discord"
	"github.com/andrewgari/starbunk-go/internal/middleware"
	"github.com/bwmarrin/discordgo"
)

var auditor = middleware.AllOf(
	middleware.NotSelf,
	middleware.HasContent,
)

func main() {
	bot.Run("RatBot", auditor, messageCreate)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "ping ratbot" {
		sender := discord.NewMessagingService(s)
		_, err := sender.SendMessage(m.ChannelID, "Pong from ratbot!")
		if err != nil {
			slog.Error("failed to send message", "bot", "ratbot", "err", err)
		}
	}
}
