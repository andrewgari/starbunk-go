# BunkBot

> Administrative backbone and YAML-driven reply bot.

## Goals & Purpose

BunkBot is the primary administrative bot for the StarBunk system. It handles
high message volume with fast reaction times and dispatches messages to a set of
configurable reply bots defined in YAML config files.

## Major Features

- **YAML-driven reply bots** — custom bots loaded from `config/replybots/*.yaml` at startup.
  Each bot has its own identity, condition logic, and response pool.
- Admin slash commands (planned).
- Webhook-based responses using `SendMessageWithIdentity`.

## YAML Config Format

Reply bots are defined in YAML files placed in the directory pointed to by the
`BUNKBOT_CONFIG_DIR` environment variable (default: `config/replybots`).

```yaml
reply-bots:
  - name: some-bot
    identity:
      type: static          # "static" | "mimic" | "random"
      botName: "Coolbot"
      avatarUrl: "https://example.com/avatar.png"
    responses:              # bot-level default response pool
      - "Default response A"
      - "Default response B"
    ignore_bots: true       # default: true — skip bot author messages
    ignore_humans: false    # default: false — process human messages
    triggers:
      - name: optional-name
        conditions:
          any_of:
            - contains_word: "hello"
            - contains_phrase: "hey there"
        responses:          # trigger-level response pool (overrides bot level)
          - "Oh hey! {start}"
```

### Identity Types

| Type | Behavior |
|------|----------|
| `static` | Fixed `botName` + `avatarUrl` — never changes |
| `mimic` | Copies a specific Discord user's name/avatar at send time. Requires `as_member: "<discord_user_id>"` |
| `random` | Picks a random non-bot guild member at send time |

Identity is resolved at message send time, so `mimic` and `random` always
reflect current Discord state.

### Condition Types (leaf)

| YAML key | Behavior |
|----------|----------|
| `contains_word: "word"` | Whole-word, case-insensitive (`\bword\b` regex) |
| `contains_phrase: "text"` | Substring, case-insensitive |
| `matches_pattern: "regex"` | Arbitrary Go regex |
| `from_user: "discord_id"` | Specific Discord user ID |
| `with_chance: 50` | Probabilistic, 0–100 (50 = 50% chance) |
| `always: true` | Always fires |

### Logical Operators

| YAML key | Behavior |
|----------|----------|
| `all_of: [...]` | AND — all conditions must match |
| `any_of: [...]` | OR — at least one must match |
| `none_of: [...]` | NOR — none may match |

Conditions nest arbitrarily deep.

### Response Templates

| Placeholder | Behavior |
|-------------|----------|
| `{start}` | First three words of the triggering message |
| `{random:2-5:e}` | Repeat char `e` between 2 and 5 times |
| `{swap_message:foo:bar}` | Replace "foo" with "bar" in the original message |

### Response Pool Priority

1. If the matching trigger defines `responses`, those are used.
2. Otherwise the bot-level `responses` pool is used.
3. If neither is set, the trigger fires silently (no message sent, warning logged).

## Architecture

- **Entry point:** `cmd/bunkbot/main.go`
- **Reply bot engine:** `internal/replybot/` (new package)
  - `config.go` — YAML struct definitions
  - `compiler.go` — `ConditionNode` → `middleware.MessageAuditor`
  - `identity.go` — YAML identity → `IdentityResolver` function
  - `response.go` — `ResponsePool` + template expansion
  - `bot.go` — `ReplyBot` runtime: `Handle` method
  - `loader.go` — `LoadFile` / `LoadDir` → `[]*ReplyBot`
- **Framework:** `internal/bot.Run` + `internal/discord.MessagingService`
- **Identity/webhook:** `internal/bot.Identity` + `DiscordIdentityProvider`
- **Config mount:** `BUNKBOT_CONFIG_DIR` env var (default: `config/replybots`)

### Message Flow

```
Discord message
  → middleware.AllOf(NotSelf, HasContent)   [global auditor in main.go]
  → for each ReplyBot in bots:
      → ReplyBot.GlobalAuditor              [ignore_bots / ignore_humans gate]
      → for each Trigger in bot:
          → Trigger.Condition.Audit()       [compiled condition tree]
          → IF match: resolve identity, pick response, SendMessageWithIdentity
          → RETURN (first match wins)
```

## Configuration in Docker

Add a volume mount to the bunkbot container:

```yaml
volumes:
  - ./config/replybots:/app/config/replybots
environment:
  - BUNKBOT_CONFIG_DIR=/app/config/replybots
```

## Edge Cases

- **Missing config directory at startup** — logged; bot stays alive but silent.
- **Invalid YAML or bad config** — load-time error logged; bot fails to start (fail-fast).
- **`mimic`/`random` identity Discord API error** — fallback to session bot identity; message still sent.
- **Webhook permission errors** — logged; no crash; no retry.
- **`ignore_bots: true` + `ignore_humans: true`** — bot loaded with warning; will never respond.

## See Also

- `cmd/bunkbot/CLAUDE.md`
- `internal/replybot/` package docs
- [[../infrastructure/Architecture|Architecture]]
