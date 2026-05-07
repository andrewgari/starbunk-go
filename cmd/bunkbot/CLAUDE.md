# BunkBot — Development Instructions

> See also: `wiki/bots/BunkBot.md`

## Goals & Purpose

Administrative backbone and general reply bot. Handles high message volume with
fast reaction times and supports webhook-based persona posting.

## Major Features

- General reply bot handlers (trigger → response mapping).
- Admin slash commands.
- Webhook-based responses via `internal/discord.MessagingService.SendMessageWithIdentity`.

## Dependencies & Architecture

- `internal/bot.Run` — event loop and session management.
- `internal/discord.MessagingService` — direct messages and webhook persona posts.
- `internal/bot.Identity` / `DiscordIdentityProvider` — resolves webhook persona from Discord.

## Edge Cases

- Webhook permission errors or rate-limit timeouts — handle gracefully, fall back to direct message.
- Race conditions on simultaneous admin commands — use locks or idempotent handlers.
- Self-message loop: always check `m.Author.ID == s.State.User.ID` before responding.
