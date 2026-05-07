# BunkBot

> Administrative backbone and general reply bot.

## Goals & Purpose

BunkBot is the primary administrative bot for the StarBunk system. It handles
high message volume with fast reaction times and can post via webhooks as custom
identities using `internal/discord.MessagingService`.

## Major Features

- General reply bot handlers (YAML-driven in JS; Go port TBD).
- Admin slash commands.
- Webhook-based responses using `SendMessageWithIdentity`.

## Dependencies & Architecture

- **Entry point:** `cmd/bunkbot/main.go`
- **Framework:** `internal/bot.Run` + `internal/discord.MessagingService`
- **Identity/webhook:** `internal/bot.Identity` + `DiscordIdentityProvider`
- Scaled for high message volume — handlers must remain lightweight.

## Edge Cases

- Webhook permission errors or timeouts.
- Race conditions on simultaneous admin commands.
- Graceful degradation when Discord API is unreachable.

## See Also

- `cmd/bunkbot/CLAUDE.md`
- [[../infrastructure/Architecture|Architecture]]
