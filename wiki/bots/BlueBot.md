# BlueBot

> Pattern-matching bot for "blue" / Blue Mage references.

## Goals & Purpose

BlueBot detects any message that references "blue" or Blue Mage and replies with
the classic catchphrase: **"Did somebody say Blu?"**

It is inspired by the starbunk-js BlueBot. The Go implementation prioritises a
clean, extensible architecture so that the trigger mechanism can be swapped or
augmented (e.g. with an LLM) without restructuring the bot.

## Architecture

### Strategy pattern

Bot-specific detection logic lives in `internal/bluebot/`. The shared
dispatcher and interface live in `internal/replybot/` — every reply-style
bot uses them. The central abstraction is the `Strategy` interface:

```go
type Strategy interface {
    Name() string
    ShouldTrigger(ctx context.Context, msg *discordgo.MessageCreate) bool
    Response(ctx context.Context, msg *discordgo.MessageCreate) string
}
```

`Bot` holds an ordered slice of strategies and dispatches to the first that
triggers. Adding behaviour means adding a `Strategy` — nothing else changes.

### Current strategies

| Strategy | Trigger | Response |
|---|---|---|
| `BlueStrategy` | Regex match on "blue" variants | `"Did somebody say Blu?"` |

### Extensibility roadmap

- **LLM trigger** — swap or layer an LLM call alongside the regex by
  implementing a new `Strategy` and passing it to `NewBot`.
- **Reply window** — a stateful strategy that opens a follow-up window after
  a blue detection, triggering on short confirmations within N minutes.
- **Enemy user** — a strategy that returns a hostile response for a specific
  user ID configured via environment variable.
- **Response variation** — extend `BlueStrategy.Response` to pick randomly
  from a list, or delegate to an LLM for generated replies.

### Pattern coverage

The regex (`bluePattern` in `internal/bluebot/blue.go`) matches:

| Variant | Example |
|---|---|
| Plain colour | `blue`, `Blue`, `BLUE` |
| Short forms | `blu`, `bluu`, `bloo` |
| Homophones | `blew` |
| French | `bleu` |
| Spanish / Portuguese | `azul` |
| German | `blau` |
| Self-reference | `bluebot` |

Word boundaries (`\b`) prevent false positives on compound words like
`bluetooth`, `blueprint`, and `blueberry`.

## Files

| File | Purpose |
|---|---|
| `cmd/bluebot/main.go` | Wires `bot.Run`, auditor, and `replybot.NewBot` |
| `internal/replybot/strategy.go` | `Strategy` interface (shared across all reply bots) |
| `internal/replybot/bot.go` | `Bot` dispatcher (shared across all reply bots) |
| `internal/bluebot/blue.go` | `BlueStrategy` — regex trigger + static response |
| `internal/bluebot/blue_test.go` | Ginkgo specs (25 cases) |

## Dependencies

- `internal/bot.Run` — event loop and session management
- `internal/discord.MessagingService` — sends replies
- No external services required (LLM integration planned)

## Edge Cases

- Compound words (`bluetooth`, `blueprint`, `blueberry`) must NOT trigger —
  enforced by word boundaries in the regex.
- The `Bot` is constructed once via `sync.Once` so stateful strategies (future)
  preserve their state across messages.
- Keep regex patterns simple and bounded to avoid ReDoS.

## See Also

- `cmd/bluebot/CLAUDE.md`
- [[../infrastructure/Architecture|Architecture]]
