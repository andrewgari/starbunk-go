# BlueBot — Development Instructions

> See also: `wiki/bots/BlueBot.md`

## Goals & Purpose

Detect any message that references "blue" or Blue Mage and reply with a witty,
character-themed response. Ported from starbunk-js BlueBot.

## Major Features

- Regex / string pattern matching across all guild messages.
- Contextual, character-specific replies.
- Optional LLM-enhanced validation to reduce false positives (not yet wired).

## Dependencies & Architecture

- `internal/bot.Run` — event loop and session management.
- `internal/discord.MessagingService` — sends replies.
- No external services required for basic pattern matching.

## Edge Cases

- Keep regex patterns simple and bounded to avoid ReDoS.
- Rate-limit replies — do not spam the channel on a burst of matching messages.
- Distinguish colloquial "blue" from intentional Blue Mage references.
