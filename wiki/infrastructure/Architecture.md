# Architecture

## Overview

Starbunk-Go is a Go monorepo. Each bot is an independent binary with its own
Docker container and Discord token. There is no shared database (contrast with
starbunk-js which uses Postgres, Redis, and Qdrant).

```
starbunk-go/
  cmd/
    bluebot/    # binary entry point + CLAUDE.md
    bunkbot/
    covabot/
    djcova/
    ratbot/
  internal/
    bot/        # bot.Run framework
    discord/    # Identity, MessageService, WebhookService
  docker/
    Dockerfile          # single multi-stage build; BOT_NAME arg selects binary
    docker-compose.yml  # local dev — builds from source
  docker-compose.yml    # production — pulls GHCR images
  .github/workflows/
    ci.yml      # PR checks
    main.yml    # build + push images on merge
    deploy.yml  # deploy to Tower server
```

## Shared Libraries

### `internal/bot`

- `bot.Run(name, auditor, handlers...)` — reads `DISCORD_TOKEN`, creates a discordgo
  session, wraps every `MessageCreate` handler with the supplied auditor, registers
  all handlers, blocks until SIGINT/SIGTERM.

### `internal/discord`

Three single-responsibility layers:

- **`Identity` / `IdentityProvider`** — first-class persona concept. `Identity{Username, Nickname, AvatarURL}` belongs to any message, not just webhook sends. `DiscordIdentityProvider` resolves live Discord identities, preferring guild-member details over global user details.

- **`WebhookService`** — internal implementation detail; callers never use it directly. Manages the full lifecycle of per-channel webhooks:
  - Lazily creates a webhook named `"Starbunk Webhook"` on first use (found by name, not by owner — all bots share one slot per channel, well within Discord's 15-webhook-per-channel limit).
  - Caches entries in a `channelID → {webhook, lastUsed}` registry.
  - Background reaper (every 1 minute) deletes webhooks idle longer than 5 minutes.
  - `Close()` stops the reaper and immediately deletes all owned webhooks for a clean shutdown.

- **`MessageService`** — the only caller-facing send API. Callers say *what* to send and *as whom*; the implementation decides how to deliver it (direct API vs webhook). `NewMessageService(s)` wires all layers internally.

  ```go
  type MessageService interface {
      // Bot's own identity — admin, errors, ephemeral messages
      SendMessage(channelID, content string) (*discordgo.Message, error)
      // Caller-provided identity — service decides transport
      SendMessageWithIdentity(channelID, content string, id Identity) (*discordgo.Message, error)
      Reply(channelID, messageID, content string) (*discordgo.Message, error)
      Edit(channelID, messageID, content string) (*discordgo.Message, error)
      Delete(channelID, messageID string) error
      // Close releases webhook resources; call on bot shutdown
      Close() error
  }
  ```

  `SendMessage` always uses the direct Discord API (no persona). `SendMessageWithIdentity` delegates to `WebhookService` — the caller never thinks about webhooks.

### `internal/llm`

- `Service` — unified abstraction for bots to interact with Large Language Models.
- Agnostic to providers (OpenAI, Anthropic, Ollama, Google).
- Explicit `ResponseSchema` allows callers to enforce format (Text, JSON, Enum) and provide validation rules or choices.

### `internal/middleware`

Composable message audit gates. Every bot must supply a `MessageAuditor` to
`bot.Run`; no `MessageCreate` handler can be invoked without passing audit.

**Interface**

```go
type MessageAuditor interface {
    Audit(s *discordgo.Session, m *discordgo.MessageCreate) bool
}
```

**Primitives** (by file)

| File | Gates |
|---|---|
| `author.go` | `NotSelf`, `NotBot`, `IsBot`, `AuthorID(id)`, `NotAuthorID(id)`, `AuthorNamed(name)`, `AuthorHasRole(roleID)` |
| `content.go` | `HasContent`, `ContentContains(substr)`, `ContentMatches(re)`, `HasAttachment` |
| `context.go` | `GuildOnly`, `DMOnly`, `InChannel(id)`, `OnWeekdays(days...)` |
| `random.go`  | `Chance(p)` — passes with probability p |

**Combinators**

```go
AllOf(auditors...)  // all must pass; short-circuits on first failure
AnyOf(auditors...)  // any must pass; short-circuits on first success
Not(auditor)        // inverts result
```

**Example composition**

```go
// Bots always fail. Non-bots pass freely, except userid 111111
// who only triggers when the message contains "bingo".
AllOf(
    NotBot,
    AnyOf(
        Not(AuthorID("111111")),
        ContentContains("bingo"),
    ),
)
```

See [[../development/MessageFiltering|Message Filtering]] for the full design.

## Bot Pattern

```go
var auditor = middleware.AllOf(
    middleware.NotSelf,
    middleware.NotBot,
    middleware.GuildOnly,
    middleware.HasContent,
)

func main() {
    bot.Run("BotName", auditor, messageCreate)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Audit has already passed. No guard needed here.
    sender := discord.NewMessageService(s)

    // Send as the bot itself (admin, errors, simple replies):
    sender.SendMessage(m.ChannelID, "response")

    // Send as a named persona (service handles transport internally):
    id := discord.Identity{Username: "CustomName", AvatarURL: "https://..."}
    sender.SendMessageWithIdentity(m.ChannelID, "response", id)
}
```

## Discord Intents

Default: `IntentsGuildMessages | IntentsMessageContent`.
DJCova additionally needs `IntentsGuildVoiceStates` — extend `bot.Run` or
create a custom initialisation.

## See Also

- [[Deployment|Deployment]]
- [[Configuration|Configuration]]
