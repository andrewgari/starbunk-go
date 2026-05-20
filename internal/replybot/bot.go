package replybot

import (
	"context"
	"log/slog"

	"github.com/andrewgari/starbunk-go/internal/discord"
	"github.com/bwmarrin/discordgo"
)

// Bot dispatches incoming Discord messages through an ordered list of
// strategies. The first strategy whose ShouldTrigger returns true wins; the
// rest are skipped. This keeps each strategy focused on a single concern and
// makes it straightforward to add, remove, or reorder behaviours.
type Bot struct {
	strategies []Strategy
	sender     discord.MessageService
}

// NewBot constructs a Bot with the given sender and strategies. Strategies
// are evaluated in the order they are passed — put higher-priority rules first.
func NewBot(sender discord.MessageService, strategies ...Strategy) *Bot {
	return &Bot{
		strategies: strategies,
		sender:     sender,
	}
}

// Handle runs each strategy in order and sends the response for the first one
// that triggers. s is the Discord session, forwarded to any ConditionedStrategy
// so conditions like AuthorHasRole can inspect guild state. ctx is forwarded to
// every strategy call so that async implementations (e.g. LLM providers) can
// respect cancellation and deadlines.
//
// If the matched strategy also implements IdentifiedStrategy and returns
// useWebhook == true, the response is sent via SendAs (webhook persona);
// otherwise it is sent as a plain direct message.
func (b *Bot) Handle(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, strategy := range b.strategies {
		if cs, ok := strategy.(ConditionedStrategy); ok {
			if !cs.Condition().Audit(s, m) {
				continue
			}
		}
		if !strategy.ShouldTrigger(ctx, m) {
			continue
		}

		resp := strategy.Response(ctx, m)

		if identified, ok := strategy.(IdentifiedStrategy); ok {
			if id, useWebhook := identified.Identity(ctx, m); useWebhook {
				if _, err := b.sender.SendAs(m.ChannelID, discord.WebhookMessage{
					Identity: id,
					Content:  resp,
				}); err != nil {
					slog.Error("replybot: failed to send webhook response",
						"strategy", strategy.Name(),
						"channel", m.ChannelID,
						"err", err,
					)
				}
				return
			}
		}

		if _, err := b.sender.Send(m.ChannelID, discord.DirectMessage{Content: resp}); err != nil {
			slog.Error("replybot: failed to send response",
				"strategy", strategy.Name(),
				"channel", m.ChannelID,
				"err", err,
			)
		}
		return
	}
}
