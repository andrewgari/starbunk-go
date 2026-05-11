// Package testenv provides a simulated Discord environment for end-to-end bot
// testing. It includes a fake discordgo.Session backed by a mock HTTP
// transport, message builder helpers, and a BotHarness that runs a full
// audit → handler → capture cycle without any real network I/O.
package testenv

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
)

// NewFakeSession returns a *discordgo.Session suitable for use in tests.
//
// The session has a populated State.User (so NotSelf auditors work correctly)
// and its HTTP client's Transport is replaced with the provided transport.
// Pass nil to use http.DefaultTransport (not recommended for unit tests).
//
// The session is never opened (dg.Open is never called), so no WebSocket
// connection is established and no Discord token is required.
func NewFakeSession(botID string, transport http.RoundTripper) *discordgo.Session {
	if transport == nil {
		transport = http.DefaultTransport
	}
	s := &discordgo.Session{
		State:       discordgo.NewState(),
		Ratelimiter: discordgo.NewRatelimiter(),
		Client:      &http.Client{Transport: transport},
	}
	s.State.User = &discordgo.User{
		ID:  botID,
		Bot: true,
	}
	return s
}
