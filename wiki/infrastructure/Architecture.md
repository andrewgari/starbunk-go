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

- **`Identity` / `IdentityProvider`** — persona model. `Identity{Username, Nickname, AvatarURL}` is required for any webhook send. `DiscordIdentityProvider` resolves live Discord identities, preferring guild-member details over global user details.

- **`WebhookService`** — manages per-channel Discord webhook lifecycle (creation, caching) and executes `WebhookMessage` payloads. `WebhookMessage.Identity` is required and validated at send time; a zero-value identity returns an error.

- **`MessageService`** — high-level send interface. `NewMessageService(s)` returns an implementation that delegates `Send` (direct API) and `SendAs` (webhook) to the appropriate layer. Implementations are fully swappable via injection.

  ```go
  type MessageService interface {
      Send(channelID string, msg DirectMessage) (*discordgo.Message, error)
      SendAs(channelID string, msg WebhookMessage) (*discordgo.Message, error)
      Reply(channelID, messageID string, msg DirectMessage) (*discordgo.Message, error)
      Edit(channelID, messageID, content string) (*discordgo.Message, error)
      Delete(channelID, messageID string) error
  }
  ```

  Message types are distinct by design — `DirectMessage` carries no identity field (compile-time enforcement), `WebhookMessage` requires one.

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
    sender.Send(m.ChannelID, discord.DirectMessage{Content: "response"})
}

// To send as a custom persona via webhook:
sender.SendAs(m.ChannelID, discord.WebhookMessage{
    Identity: discord.Identity{Username: "CustomName", AvatarURL: "https://..."},
    Content:  "response",
})
```

## Discord Intents

Default: `IntentsGuildMessages | IntentsMessageContent`.
DJCova additionally needs `IntentsGuildVoiceStates` — extend `bot.Run` or
create a custom initialisation.

## See Also

- [[Deployment|Deployment]]
- [[Configuration|Configuration]]
