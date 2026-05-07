# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run all tests
go test ./...

# Run tests in a specific package
go test ./internal/...

# Run a single test (Ginkgo)
go test ./internal/... -run "TestInternal/Config"

# Build a specific bot
go build -o bot ./cmd/bunkbot

# Build all bots
for bot in bluebot bunkbot covabot djcova ratbot; do go build -o bin/$bot ./cmd/$bot; done

# Run a bot locally (requires DISCORD_TOKEN env var)
DISCORD_TOKEN=<token> go run ./cmd/bunkbot

# Build and run all containers
docker compose -f docker/docker-compose.yml up -d --build

# Build a single container
docker compose -f docker/docker-compose.yml up -d --build bunkbot
```

## Architecture

This is a **Go monorepo** housing 5 independent Discord bots (`bluebot`, `bunkbot`, `covabot`, `djcova`, `ratbot`), each with its own binary entry point in `cmd/<botname>/main.go` and its own Discord token.

### Shared internal libraries (`internal/`)

- **`internal/bot`** — Core bot framework.
  - `bot.Run(name, handlers...)` reads `DISCORD_TOKEN` from env, creates a `discordgo.Session`, registers event handlers, and blocks until SIGINT/SIGTERM.
  - `Identity` / `IdentityProvider` — Persona model used when a bot needs to impersonate a user via webhooks. `DiscordIdentityProvider` resolves a user's identity from Discord (prefers guild-member details over global user details).

- **`internal/discord`** — Messaging abstraction.
  - `MessagingService` interface wraps `discordgo.Session` for sending, replying, editing, and deleting messages. `SendMessageWithIdentity` uses a per-channel webhook (created lazily) to post as a custom user/avatar.

### Bot pattern

Every bot follows the same pattern:

```go
func main() {
    bot.Run("BotName", messageCreate)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID { return }
    // respond to m.Content
    sender := discord.NewMessagingService(s)
    sender.SendMessage(m.ChannelID, "response")
}
```

### Discord intents

`bot.Run` sets `IntentsGuildMessages | IntentsMessageContent`. If a new bot needs additional intents (voice, reactions, etc.) it must extend `bot.Run` or create a custom initialization.

### Environment variables

| Variable | Purpose |
|---|---|
| `DISCORD_TOKEN` | Token used by `bot.Run` at runtime |
| `STARBUNK_TOKEN` | Fallback token in Docker Compose |
| `{BOTNAME}_TOKEN` | Per-bot override (e.g. `BUNKBOT_TOKEN`) |
| `CLOUD_LLM_PROVIDER` / `CLOUD_LLM_API_KEY` | Cloud LLM (Gemini) — not yet wired |
| `LOCAL_LLM_PROVIDER` / `LOCAL_LLM_API_KEY` | Local LLM (Ollama) — not yet wired |

In Docker Compose, each service resolves its token as `${BOTNAME_TOKEN:-${STARBUNK_TOKEN}}`.

### Testing

Tests use the **Ginkgo v2 / Gomega** BDD framework. Test files use the `_test` package suffix. To add tests for a new package, create a suite bootstrap:

```go
func TestFoo(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Foo Suite")
}
```

### Adding a new bot

1. Create `cmd/<newbot>/main.go` calling `bot.Run`.
2. Add a service entry in `docker/docker-compose.yml` with `BOT_NAME: <newbot>` and `DISCORD_TOKEN: ${NEWBOT_TOKEN:-${STARBUNK_TOKEN}}`.

## Wiki maintenance

Future agents should keep this CLAUDE.md updated whenever new bots, shared packages, or significant architectural patterns are introduced.
