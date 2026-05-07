# RatBot

> Rat-trigger response bot.

## Goals & Purpose

RatBot monitors guild messages for rat-related triggers and fires a response.
Specific trigger patterns and response content are TBD — this bot exists in the
Go monorepo but has no equivalent in starbunk-js.

## Major Features

- Pattern matching for rat-related keywords or phrases.
- Configurable response content.

## Dependencies & Architecture

- **Entry point:** `cmd/ratbot/main.go`
- **Framework:** `internal/bot.Run` + `internal/discord.MessagingService`
- Lightweight; same pattern as BlueBot.

## Edge Cases

- Avoid false positives on unrelated messages containing substrings like "rat".
- Rate-limit replies to prevent spam.

## See Also

- `cmd/ratbot/CLAUDE.md`
- [[BlueBot|BlueBot]] — same pattern-matching architecture
