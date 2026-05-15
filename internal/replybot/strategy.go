// Package replybot provides the shared Strategy interface and Bot dispatcher
// used by every reply-style bot in starbunk-go.
//
// A Strategy encapsulates one trigger condition and one response. A Bot holds
// an ordered list of strategies and dispatches each incoming message to the
// first strategy that fires. New behaviours are added by implementing Strategy
// and passing it to NewBot — the dispatcher never changes.
//
// Today strategies are regex-based; the interface is forward-compatible with
// async implementations such as LLM calls because ctx is threaded through
// every method.
package replybot

import (
	"context"

	"github.com/andrewgari/starbunk-go/internal/middleware"
	"github.com/bwmarrin/discordgo"
)

// Strategy is the single extensibility seam for reply bots. Implement this
// interface to add a trigger mechanism (regex, keyword list, LLM, …) or a
// response style (static, randomised, LLM-generated, …) without touching
// the Bot dispatcher or any other strategy.
//
// ctx is provided so that future async implementations (e.g. LLM providers
// with cancellation deadlines) can respect caller context.
type Strategy interface {
	// Name identifies this strategy in logs and future metrics.
	Name() string

	// ShouldTrigger returns true when this strategy wants to respond to msg.
	ShouldTrigger(ctx context.Context, msg *discordgo.MessageCreate) bool

	// Response returns the text to send. Only called after ShouldTrigger
	// returns true, so implementations may assume the message matched.
	Response(ctx context.Context, msg *discordgo.MessageCreate) string
}

// ConditionedStrategy is an optional extension of Strategy. When a strategy
// implements this interface, Bot.Handle evaluates the condition before calling
// ShouldTrigger. Messages that fail the condition are silently skipped for
// this strategy; the remaining strategies continue to be evaluated.
//
// Use WithCondition to compose any existing Strategy with a condition without
// modifying the strategy struct itself.
type ConditionedStrategy interface {
	Strategy
	Condition() middleware.MessageAuditor
}
