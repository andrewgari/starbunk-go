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

- `bot.Run(name, handlers...)` — reads `DISCORD_TOKEN`, creates a discordgo session,
  registers handlers, blocks until SIGINT/SIGTERM.
- `Identity` / `IdentityProvider` — persona model for webhook impersonation.
  `DiscordIdentityProvider` prefers guild-member details over global user details.

### `internal/discord`

- `MessagingService` — interface over discordgo for send, reply, edit, delete.
- `SendMessageWithIdentity` — creates/reuses a per-channel webhook to post as a
  custom user/avatar.

## Bot Pattern

```go
func main() {
    bot.Run("BotName", messageCreate)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID { return }
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
