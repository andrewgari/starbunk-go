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

- Keep regex patterns simple and bounded to avoid ReDoS (Regex Denial of
  Service). Never use unbounded quantifiers on large character classes.
- **Rate-limit replies** — use separate windows per reply type:
  - Standard pattern matches: ~5-minute cooldown per channel.
  - Rare / enemy-themed responses: ~24-hour cooldown per channel.
- Distinguish colloquial "blue" (e.g. colour, mood) from intentional Blue Mage
  references. Optional LLM validation can reduce false positives once wired.
- Avoid triggering on messages sent by other bots (check `m.Author.Bot`).
