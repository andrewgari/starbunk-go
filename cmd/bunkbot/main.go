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
	bot.Run("BunkBot", auditor, messageCreate)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "ping bunkbot" {
		sender := discord.NewMessageService(s)
		if _, err := sender.Send(m.ChannelID, discord.DirectMessage{Content: "Pong from bunkbot!"}); err != nil {
			slog.Error("failed to send message", "bot", "bunkbot", "err", err)
		}
	}
}
