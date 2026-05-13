package replybot

import (
	"context"
	"log"

	"github.com/andrewgari/starbunk-go/internal/discord"
	"github.com/bwmarrin/discordgo"
)

// Bot dispatches incoming Discord messages through an ordered list of
// strategies. The first strategy whose ShouldTrigger returns true wins; the
// rest are skipped. This keeps each strategy focused on a single concern and
// makes it straightforward to add, remove, or reorder behaviours.
type Bot struct {
	strategies []Strategy
	sender     discord.MessagingService
}

// NewBot constructs a Bot with the given sender and strategies. Strategies
// are evaluated in the order they are passed — put higher-priority rules first.
func NewBot(sender discord.MessagingService, strategies ...Strategy) *Bot {
	return &Bot{
		strategies: strategies,
		sender:     sender,
	}
}

// Handle is the discordgo MessageCreate event handler. It runs each strategy
// in order and sends the response for the first one that triggers.
func (b *Bot) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	ctx := context.Background()
	for _, strategy := range b.strategies {
		if strategy.ShouldTrigger(ctx, m) {
			resp := strategy.Response(ctx, m)
			if _, err := b.sender.SendMessage(m.ChannelID, resp); err != nil {
				log.Printf("[replybot] %s: failed to send response: %v", strategy.Name(), err)
			}
			return
		}
	}
}
