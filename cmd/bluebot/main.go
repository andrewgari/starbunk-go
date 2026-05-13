package main

import (
	"sync"

	"github.com/andrewgari/starbunk-go/internal/bluebot"
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

var (
	blueBotOnce sync.Once
	blueBot     *bluebot.Bot
)

func main() {
	bot.Run("BlueBot", auditor, messageCreate)
}

// messageCreate is the registered Discord event handler. The Bot is
// constructed once on the first message so it can hold state across calls
// (e.g. reply-window timers or an LLM client added in future).
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	blueBotOnce.Do(func() {
		blueBot = bluebot.NewBot(
			discord.NewMessagingService(s),
			bluebot.BlueStrategy{},
		)
	})
	blueBot.Handle(s, m)
}
