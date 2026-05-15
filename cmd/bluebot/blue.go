package main

import (
	"context"
	"regexp"

	"github.com/andrewgari/starbunk-go/internal/replybot"
	"github.com/bwmarrin/discordgo"
)

// Compile-time assertion: BlueStrategy must satisfy the shared Strategy interface.
var _ replybot.Strategy = BlueStrategy{}

// bluePattern matches any plausible reference to "blue" — the colour, the
// job, the mage, common homophones, and a handful of other-language spellings.
//
// Word boundaries (\b) prevent false positives on compound words like
// "bluetooth", "blueprint", or "blueberry".
var bluePattern = regexp.MustCompile(`(?i)\b(bluebot|bloo+|bleu|blew|azul|blau|blu+|blue?)\b`)

// BlueStrategy triggers on any message that contains a recognisable reference
// to "blue" and replies with the classic catchphrase.
//
// It is intentionally stateless. Strategies that need state (reply windows,
// cooldown timers, enemy-user tracking) should be separate replybot.Strategy
// implementations composed alongside this one via replybot.NewBot.
type BlueStrategy struct{}

func (BlueStrategy) Name() string { return "BlueStrategy" }

func (BlueStrategy) ShouldTrigger(_ context.Context, msg *discordgo.MessageCreate) bool {
	return bluePattern.MatchString(msg.Content)
}

func (BlueStrategy) Response(_ context.Context, _ *discordgo.MessageCreate) string {
	return "Did somebody say Blu?"
}
