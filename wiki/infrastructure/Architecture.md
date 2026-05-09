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
    bot/        # bot.Run, Identity, IdentityProvider
    discord/    # MessagingService interface + implementation
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
- `Identity` / `IdentityProvider` — persona model for webhook impersonation.
  `DiscordIdentityProvider` prefers guild-member details over global user details.

### `internal/discord`

- `MessagingService` — interface over discordgo for send, reply, edit, delete.
- `SendMessageWithIdentity` — creates/reuses a per-channel webhook to post as a
  custom user/avatar.

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
    sender := discord.NewMessagingService(s)
    sender.SendMessage(m.ChannelID, "response")
}
```

## Discord Intents

Default: `IntentsGuildMessages | IntentsMessageContent`.
DJCova additionally needs `IntentsGuildVoiceStates` — extend `bot.Run` or
create a custom initialisation.

## See Also

- [[Deployment|Deployment]]
- [[Configuration|Configuration]]
