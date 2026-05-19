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
	middleware.NotBot,
	middleware.GuildOnly,
	middleware.HasContent,
)

func main() {
	bot.Run("CovaBot", auditor, messageCreate)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "ping covabot" {
		sender := discord.NewMessagingService(s)
		_, err := sender.SendMessage(m.ChannelID, "Pong from covabot!")
		if err != nil {
			slog.Error("failed to send message", "bot", "covabot", "err", err)
		}
	}
}
