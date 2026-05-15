package replybot

import (
	"context"

	"github.com/andrewgari/starbunk-go/internal/middleware"
	"github.com/bwmarrin/discordgo"
)

// WithCondition wraps any Strategy with a MessageAuditor condition.
// Bot.Handle checks the condition before calling ShouldTrigger, so the inner
// strategy is only reached when the condition passes. This lets you compose
// conditions at the call site (e.g. in NewBot) without modifying the strategy.
//
// Example — BunkBot with a bot-only and a human-only strategy:
//
//	replybot.NewBot(sender,
//	    replybot.WithCondition(middleware.IsBot,  botBotStrategy),
//	    replybot.WithCondition(middleware.NotBot, humanOnlyStrategy),
//	)
func WithCondition(cond middleware.MessageAuditor, s Strategy) Strategy {
	return conditionedStrategy{cond: cond, inner: s}
}

type conditionedStrategy struct {
	cond  middleware.MessageAuditor
	inner Strategy
}

func (c conditionedStrategy) Name() string { return c.inner.Name() }
func (c conditionedStrategy) Condition() middleware.MessageAuditor { return c.cond }

func (c conditionedStrategy) ShouldTrigger(ctx context.Context, msg *discordgo.MessageCreate) bool {
	return c.inner.ShouldTrigger(ctx, msg)
}

func (c conditionedStrategy) Response(ctx context.Context, msg *discordgo.MessageCreate) string {
	return c.inner.Response(ctx, msg)
}
