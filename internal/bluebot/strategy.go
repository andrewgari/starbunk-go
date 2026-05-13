// Package bluebot implements the BlueBot strategy engine.
//
// The central abstraction is Strategy — anything that can decide whether to
// respond to a message and what to say. Today that is a regex; in future it
// can be an LLM call, a chain of sub-strategies, or anything else that
// satisfies the interface without changing the Bot or the caller.
package bluebot

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// Strategy is the extensibility seam for BlueBot. Implement this interface
// to add a new trigger mechanism (regex, LLM, keyword list, …) or a new
// response style (static, randomised, LLM-generated, …) without touching
// the Bot dispatcher.
//
// ctx is threaded through every call so that future async implementations
// (e.g. an LLM provider with a deadline) can respect cancellation.
type Strategy interface {
	// Name identifies this strategy in logs and future metrics.
	Name() string

	// ShouldTrigger returns true when this strategy wants to respond to msg.
	ShouldTrigger(ctx context.Context, msg *discordgo.MessageCreate) bool

	// Response returns the text to send. Only called after ShouldTrigger
	// returns true, so implementations may assume the message matched.
	Response(ctx context.Context, msg *discordgo.MessageCreate) string
}
