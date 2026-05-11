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

## E2E Testing

Tests live in `cmd/bluebot/e2e_test.go` (`package main`), run via `go test ./cmd/bluebot/...`.

**Auditor policy under test:** `AllOf(NotSelf, NotBot, GuildOnly, HasContent)`

| Scenario | Expected |
|---|---|
| Human guild message with content | Audit passes |
| Self message | Blocked by `NotSelf` |
| Bot author | Blocked by `NotBot` |
| DM (no GuildID) | Blocked by `GuildOnly` |
| Empty content | Blocked by `HasContent` |
| `"ping bluebot"` | Audit passes, sends `"Pong from bluebot!"` |
| Unrecognised content | Audit passes, no reply |

On test failure, `testenv.FormatScenarioFailure` automatically appends the full audit tree (each check: PASSED / FAILED / IGNORED) to the Ginkgo failure message.

## See Also

- `cmd/bluebot/CLAUDE.md`
- [[../infrastructure/Architecture|Architecture]]
