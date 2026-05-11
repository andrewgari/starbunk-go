package replybot

import (
	"log"

	"github.com/andrewgari/starbunk-go/internal/discord"
	"github.com/andrewgari/starbunk-go/internal/middleware"
	"github.com/bwmarrin/discordgo"
)

// CompiledTrigger is a trigger after its condition has been compiled into
// a live MessageAuditor. Pool may be nil, in which case the bot-level pool is used.
type CompiledTrigger struct {
	Name      string
	Condition middleware.MessageAuditor
	Pool      *ResponsePool // nil = fall through to bot-level pool
}

// ReplyBot is the runtime representation of one YAML bot entry.
// It holds compiled triggers and responds to the first matching trigger.
type ReplyBot struct {
	Name             string
	GlobalAuditor    middleware.MessageAuditor
	IdentityResolver IdentityResolver
	Triggers         []CompiledTrigger
	BotPool          *ResponsePool // may be nil if no bot-level responses
	Messaging        discord.MessagingService
}

// Handle evaluates the message against the global auditor and all triggers.
// The first matching trigger fires: identity is resolved, a response is picked,
// and the message is sent. No further triggers are evaluated after a match.
func (b *ReplyBot) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !b.GlobalAuditor.Audit(s, m) {
		return
	}

	for _, trigger := range b.Triggers {
		if !trigger.Condition.Audit(s, m) {
			continue
		}

		pool := trigger.Pool
		if pool == nil {
			pool = b.BotPool
		}
		if pool == nil {
			log.Printf("replybot %q: trigger %q matched but no responses configured", b.Name, trigger.Name)
			return
		}

		id, err := b.IdentityResolver(s, m)
		if err != nil {
			log.Printf("replybot %q: identity resolution failed, using session fallback: %v", b.Name, err)
			id = id.Resolve(s)
		}

		content := pool.Pick(m.Content)
		if _, err := b.Messaging.SendMessageWithIdentity(m.ChannelID, content, id.Username, id.AvatarURL); err != nil {
			log.Printf("replybot %q: failed to send message: %v", b.Name, err)
		}
		return
	}
}
