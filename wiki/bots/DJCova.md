# DJCova

> Voice channel music streaming service.

## Goals & Purpose

DJCova joins Discord voice channels on demand and streams YouTube audio. It
manages a per-guild queue so each server has independent playback state. Ported
from starbunk-js DJCova.

## Major Features

- `/play <youtube-url>` command — joins voice and streams audio.
- Per-guild queue management (add, skip, clear).
- Voice channel state management (join, leave, reconnect).

## Dependencies & Architecture

- **Entry point:** `cmd/djcova/main.go`
- **Framework:** `internal/bot.Run` — **requires voice intents** (not in default `bot.Run`; needs extension)
- **Audio:** ffmpeg + discordgo voice (not yet wired in Go port)
- CPU-intensive; must not block the event loop.

## Edge Cases

- Voice connection health monitoring and reconnection.
- Concurrent `/play` requests and queue races.
- YouTube playback errors or geo-restricted videos.
- Proper cleanup of ffmpeg processes on disconnect or crash.

## E2E Testing

Tests live in `cmd/djcova/e2e_test.go` (`package main`), run via `go test ./cmd/djcova/...`.

**Auditor policy under test:** `AllOf(NotSelf, HasContent)`

| Scenario | Expected |
|---|---|
| Human guild message with content | Audit passes |
| Self message | Blocked by `NotSelf` |
| Empty content | Blocked by `HasContent` |
| Bot author (intentional) | Audit passes — DJCova allows bot-authored messages |
| DM (intentional) | Audit passes — DJCova allows DMs |
| `"ping djcova"` | Audit passes, sends `"Pong from djcova!"` |
| Unrecognised content | Audit passes, no reply |

## See Also

- `cmd/djcova/CLAUDE.md`
- [[../infrastructure/Architecture|Architecture]] — note on extending `bot.Run` for voice intents
