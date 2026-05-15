package main

import (
	"context"
	"sync"

	"github.com/andrewgari/starbunk-go/internal/bot"
	"github.com/andrewgari/starbunk-go/internal/discord"
	"github.com/andrewgari/starbunk-go/internal/middleware"
	"github.com/andrewgari/starbunk-go/internal/replybot"
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
	blueBot     *replybot.Bot
)

func main() {
	bot.Run("BlueBot", auditor, messageCreate)
}

// messageCreate is the registered Discord event handler. The Bot is
// constructed once on the first message so it can hold state across calls
// (e.g. reply-window timers or an LLM client added in future).
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	blueBotOnce.Do(func() {
		blueBot = replybot.NewBot(
			discord.NewMessagingService(s),
			BlueStrategy{},
		)
	})
	blueBot.Handle(context.Background(), s, m)
}
