package main

import (
	"log"
	"os"
	"sync"

	"github.com/andrewgari/starbunk-go/internal/bot"
	"github.com/andrewgari/starbunk-go/internal/discord"
	"github.com/andrewgari/starbunk-go/internal/middleware"
	"github.com/andrewgari/starbunk-go/internal/replybot"
	"github.com/bwmarrin/discordgo"
)

var auditor = middleware.AllOf(
	middleware.NotSelf,
	middleware.HasContent,
)

func main() {
	bot.Run("BunkBot", auditor, buildHandler())
}

// buildHandler returns a handler that lazily loads YAML reply-bot configs on
// the first message. The sync.Once ensures thread-safe single-init under
// discordgo's concurrent message delivery.
func buildHandler() func(*discordgo.Session, *discordgo.MessageCreate) {
	configDir := os.Getenv("BUNKBOT_CONFIG_DIR")
	if configDir == "" {
		configDir = "config/replybots"
	}

	var (
		once sync.Once
		bots []*replybot.ReplyBot
	)

	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		once.Do(func() {
			var err error
			bots, err = replybot.LoadDir(configDir, discord.NewMessagingService(s))
			if err != nil {
				log.Printf("bunkbot: failed to load reply-bots from %s: %v", configDir, err)
			}
			log.Printf("bunkbot: loaded %d reply-bot(s) from %s", len(bots), configDir)
		})
		for _, b := range bots {
			b.Handle(s, m)
		}
	}
}
