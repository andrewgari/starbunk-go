package testenv

import (
	"github.com/andrewgari/starbunk-go/internal/middleware"
	"github.com/bwmarrin/discordgo"
)

// BotHarness wires together a fake Discord session, a transport interceptor,
// an auditor, and a message handler for end-to-end scenario testing.
//
// Use NewBotHarness to construct one, then call Run for each test scenario.
type BotHarness struct {
	// Session is the fake discordgo.Session used for both auditing and handler
	// invocation. Exposed so tests can add guild/member state if needed.
	Session   *discordgo.Session
	// Transport captures all outbound Discord API calls the handler makes.
	Transport *FakeDiscordTransport

	auditor middleware.MessageAuditor
	handler func(*discordgo.Session, *discordgo.MessageCreate)
}

// ScenarioResult holds the full outcome of one BotHarness.Run call.
type ScenarioResult struct {
	// AuditPassed is true if the auditor allowed the message through.
	AuditPassed bool
	// AuditTrace is the full auditor tree. Always populated regardless of outcome.
	AuditTrace middleware.AuditNode
	// Messages contains all messages the handler attempted to send.
	// Empty if AuditPassed is false or the handler sent nothing.
	Messages []CapturedMessage
}

// NewBotHarness constructs a BotHarness for the given bot identity.
// botID is used as the session's State.User.ID (for NotSelf checks).
func NewBotHarness(
	botID string,
	auditor middleware.MessageAuditor,
	handler func(*discordgo.Session, *discordgo.MessageCreate),
) *BotHarness {
	transport := &FakeDiscordTransport{}
	session := NewFakeSession(botID, transport)
	return &BotHarness{
		Session:   session,
		Transport: transport,
		auditor:   auditor,
		handler:   handler,
	}
}

// Run simulates one incoming Discord message through the full bot pipeline:
//  1. Runs TraceAudit on the auditor.
//  2. If the auditor passes, invokes the handler synchronously.
//  3. Returns the full ScenarioResult including the audit tree and any
//     messages the handler attempted to send.
//
// The Transport is reset before each call so captured messages belong only
// to the current scenario.
func (h *BotHarness) Run(msg *discordgo.MessageCreate) ScenarioResult {
	h.Transport.Reset()
	passed, trace := middleware.TraceAudit(h.auditor, h.Session, msg)
	result := ScenarioResult{
		AuditPassed: passed,
		AuditTrace:  trace,
	}
	if passed {
		h.handler(h.Session, msg)
		result.Messages = h.Transport.Messages()
	}
	return result
}
