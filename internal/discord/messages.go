package discord

// DirectMessage is sent by the bot under its own identity.
// It carries no persona — use WebhookMessage when a custom display name or
// avatar is required.
type DirectMessage struct {
	Content string
}

// WebhookMessage is sent via a Discord webhook so it can appear under a custom
// display name and avatar. Identity is required: a zero-value Identity will be
// rejected at send time (see Identity.IsValid).
type WebhookMessage struct {
	Identity Identity // required
	Content  string
}
