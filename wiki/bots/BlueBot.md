# BlueBot

> Pattern-matching bot for "blue" / Blue Mage references.

## Goals & Purpose

BlueBot detects any message that references "blue" or Blue Mage and replies with
a witty, character-themed response. It is a port of the starbunk-js BlueBot.

## Major Features

- Regex / string pattern matching across all guild messages.
- Contextual, character-specific replies.
- Optional LLM-enhanced validation to avoid false positives (not yet wired in Go port).

## Dependencies & Architecture

- **Entry point:** `cmd/bluebot/main.go`
- **Framework:** `internal/bot.Run` + `internal/discord.MessagingService`
- **External:** none (LLM integration planned)
- Lightweight footprint — optimised for regex speed.

## Edge Cases

- Avoid ReDoS: keep patterns simple and bounded.
- Rate-limit replies to prevent channel spam.
- Distinguish colloquial "blue" from intentional Blue Mage references.

## See Also

- `cmd/bluebot/CLAUDE.md`
- [[../infrastructure/Architecture|Architecture]]
