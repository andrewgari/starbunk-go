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

# Build and run all containers (local dev — builds from source)
docker compose -f docker/docker-compose.yml up -d --build

# Build a single container (local dev)
docker compose -f docker/docker-compose.yml up -d --build bunkbot

# Validate DevOps file consistency (REQUIRED after any bot or CI/CD change)
bash scripts/devops-validate.sh
```

---

## !! DevOps File Maintenance — MANDATORY !!

> **This section applies to every agent and every task. Skipping it is not
> acceptable.** The CI pipeline enforces this check — drift will cause the
> `validate_devops` job to fail and block the entire pipeline.

### The rule

Every bot that lives under `cmd/<botname>/` **must** be registered in **all
six** of the following files. They must always be kept in sync with each other:

| File | What to update |
|---|---|
| `docker-compose.yml` | Add a service with `image: ghcr.io/andrewgari/starbunk-go-<bot>:${IMAGE_TAG:-latest}` |
| `docker/docker-compose.yml` | Add a service with `BOT_NAME: <bot>` build arg |
| `.github/workflows/ci.yml` | Add `<bot>` to the `build` job matrix |
| `.github/workflows/main.yml` | Add `cmd/<bot>/**` to the `paths-filter` block and to the `workflow_dispatch` matrix |
| `scripts/deployment/health-check.sh` | Add `"<bot>"` to the `EXPECTED_SERVICES` array |
| `CLAUDE.md` | Update the bot list everywhere it appears in this file |

### Validation step — run this after every relevant change

After **any** task that adds, removes, or renames a bot, or that touches any
of the DevOps files above, you **must** run:

```bash
bash scripts/devops-validate.sh
```

Fix every `FAIL` line before marking the task complete. This script is also
run as the `validate_devops` job at the start of every CI and CD pipeline, so
any drift will block merges.

### When does this apply?

Run the validation after any task involving:

- Adding a new bot (`cmd/<newbot>/`)
- Removing or renaming a bot
- Editing `docker-compose.yml`, `docker/docker-compose.yml`, or `docker/Dockerfile`
- Editing any file under `.github/workflows/`
- Editing `scripts/deployment/health-check.sh`
- Editing this `CLAUDE.md` (keep the bot lists here in sync too)

---

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

### Docker Compose files

There are **two** compose files with different purposes — do not confuse them:

| File | Purpose |
|---|---|
| `docker-compose.yml` | **Production** — pulls pre-built images from GHCR. This is the file deployed to the Tower server by `deploy.yml`. |
| `docker/docker-compose.yml` | **Local dev** — builds images from source using `docker/Dockerfile` with a `BOT_NAME` build arg. |

The production `docker-compose.yml` requires a `stack.env` file on the server containing the bot tokens. `IMAGE_TAG` is injected by `deploy.sh` (defaults to `latest`).

### CI/CD pipeline

| Workflow | Trigger | What it does |
|---|---|---|
| `ci.yml` | PRs to `main` | Validate DevOps consistency, go vet, go test, build each bot binary |
| `main.yml` | Merge to `main` | Validate DevOps consistency, test, build+push Docker images to GHCR, tag `:latest`, create git build tag |
| `deploy.yml` | After `main.yml` succeeds | Tailscale SSH to Tower, stage new compose file, pull images, restart services, health check |

GHCR image names follow the pattern `ghcr.io/andrewgari/starbunk-go-<bot>:<tag>`.

### Testing

Tests use the **Ginkgo v2 / Gomega** BDD framework. Test files use the `_test` package suffix. To add tests for a new package, create a suite bootstrap:

```go
func TestFoo(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Foo Suite")
}
```

---

## Adding a new bot — complete checklist

> After completing every step, run `bash scripts/devops-validate.sh`.
> All checks must pass before the work is done.

1. **Create** `cmd/<newbot>/main.go` calling `bot.Run`.

2. **`docker-compose.yml`** — add a service block:
   ```yaml
   <newbot>:
     image: ghcr.io/andrewgari/starbunk-go-<newbot>:${IMAGE_TAG:-latest}
     container_name: starbunk-go-<newbot>
     restart: unless-stopped
     environment:
       - DISCORD_TOKEN=${NEWBOT_TOKEN:-${STARBUNK_TOKEN}}
     logging:
       driver: "json-file"
       options:
         max-size: "10m"
         max-file: "3"
     labels:
       - "com.centurylinklabs.watchtower.enable=true"
   ```

3. **`docker/docker-compose.yml`** — add a service block:
   ```yaml
   <newbot>:
     build:
       context: ..
       dockerfile: docker/Dockerfile
       args:
         BOT_NAME: <newbot>
     container_name: starbunk-go-<newbot>
     restart: unless-stopped
     environment:
       - DISCORD_TOKEN=${NEWBOT_TOKEN:-${STARBUNK_TOKEN}}
     logging:
       driver: "json-file"
       options:
         max-size: "10m"
         max-file: "3"
   ```

4. **`.github/workflows/ci.yml`** — add `<newbot>` to the build matrix:
   ```yaml
   matrix:
     bot: [bluebot, bunkbot, covabot, djcova, ratbot, <newbot>]
   ```

5. **`.github/workflows/main.yml`** — add two entries:
   - In the `paths-filter` block:
     ```yaml
     <newbot>:
       - 'cmd/<newbot>/**'
     ```
   - In the `workflow_dispatch` matrix line and `BOTS=()` fallback:
     ```bash
     echo 'bots=["bluebot","bunkbot","covabot","djcova","ratbot","<newbot>"]' >> $GITHUB_OUTPUT
     ```

6. **`scripts/deployment/health-check.sh`** — add `"<newbot>"` to `EXPECTED_SERVICES`:
   ```bash
   EXPECTED_SERVICES=(bluebot bunkbot covabot djcova ratbot <newbot>)
   ```

7. **`CLAUDE.md`** — update this file: the bot list in the Architecture section and this checklist.

8. **Run validation**:
   ```bash
   bash scripts/devops-validate.sh
   ```
   Fix every `FAIL` before committing.

---

## Wiki maintenance

Future agents must keep this CLAUDE.md updated whenever new bots, shared packages, CI/CD changes, or significant architectural patterns are introduced. The DevOps validation section above is authoritative — update it if the set of required files changes.
