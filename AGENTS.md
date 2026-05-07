# AGENTS.md

Canonical agent guide for `starbunk-go`. All AI coding tools should read this
file. Claude Code loads it automatically via `CLAUDE.md`; other tools should
read it directly.

---

## !! MANDATORY: KEEP THE WIKI UP TO DATE !!

> **This rule applies to every agent, every task, without exception.**
> There is no situation where skipping a wiki update is acceptable.

1. **Before starting any task** — read the relevant `wiki/` page(s) for the area you will touch.
2. **After completing any task** — update the relevant wiki page(s) to reflect any changes to architecture, behavior, configuration, or patterns.
3. **For every significant change or PR** — add an entry to `wiki/Changelog.md` under today's date.
4. If a wiki page does not exist for the area you are working in, **create it**.

The wiki lives at `wiki/` in the repo root. Start at `wiki/Home.md`.

---

## !! MANDATORY: USE AVAILABLE SKILLS AND TOOLS PROACTIVELY !!

> **Do not wait to be told.** When a situation matches an available skill or
> tool capability, use it immediately without prompting.

Examples:
- Starting a non-trivial implementation → use a plan / task breakdown first
- Code has been written or changed → review it for simplicity and quality
- A PR needs deployment → use the deploy skill/tool
- Tests are available → run them before declaring a task done

If your agent environment exposes named skills or slash commands, invoke them
when the task matches — don't describe what you would do, just do it.

---

## Definition of Done

A task is **not complete** until:

1. All CI checks pass (`Validate DevOps Consistency`, `Lint`, `Test`).
2. If a PR was opened — it has at least one approval and all checks are green.
3. `bash scripts/devops-validate.sh` exits cleanly (if any bot or CI/CD file was touched).
4. The relevant `wiki/` page(s) have been updated.
5. An entry has been added to `wiki/Changelog.md`.

"The code works locally" is not done. "The PR is open" is not done.

---

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
| `AGENTS.md` | Update the bot list everywhere it appears in this file |

### Validation step — run this after every relevant change

After **any** task that adds, removes, or renames a bot, or that touches any
of the DevOps files above, you **must** run:

```bash
bash scripts/devops-validate.sh
```

Fix every `FAIL` line before marking the task complete.

### When does this apply?

Run the validation after any task involving:

- Adding a new bot (`cmd/<newbot>/`)
- Removing or renaming a bot
- Editing `docker-compose.yml`, `docker/docker-compose.yml`, or `docker/Dockerfile`
- Editing any file under `.github/workflows/`
- Editing `scripts/deployment/health-check.sh`
- Editing `AGENTS.md` (keep the bot lists here in sync too)

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

| File | Purpose |
|---|---|
| `docker-compose.yml` | **Production** — pulls pre-built images from GHCR |
| `docker/docker-compose.yml` | **Local dev** — builds images from source using `docker/Dockerfile` with a `BOT_NAME` build arg |

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

## Bots

Each bot has a `cmd/<bot>/CLAUDE.md` with its goals, features, and edge cases.
Per-bot wiki pages live at `wiki/bots/<Bot>.md`.

Current bots: `bluebot`, `bunkbot`, `covabot`, `djcova`, `ratbot`.

### bluebot
Pattern-matching bot. Detects references to "blue" or Blue Mage in messages and
replies with contextual or character-themed responses. Ported from starbunk-js.
See `cmd/bluebot/CLAUDE.md` and `wiki/bots/BlueBot.md`.

### bunkbot
Administrative backbone and general reply bot. Handles high message volume with
fast reaction times. May use webhooks to post as other identities. Ported from
starbunk-js. See `cmd/bunkbot/CLAUDE.md` and `wiki/bots/BunkBot.md`.

### covabot
AI personality emulator. Responds to conversational mentions with LLM-driven
replies that mimic a specific user's tone. Depends on an LLM provider (Ollama /
Gemini / OpenAI). Ported from starbunk-js. See `cmd/covabot/CLAUDE.md` and
`wiki/bots/CovaBot.md`.

### djcova
Voice channel music streaming service. Joins voice on demand, plays YouTube
audio, manages a per-guild queue. Requires additional voice intents. Ported from
starbunk-js. See `cmd/djcova/CLAUDE.md` and `wiki/bots/DJCova.md`.

### ratbot
Watches for rat-related triggers in messages and responds accordingly. See
`cmd/ratbot/CLAUDE.md` and `wiki/bots/RatBot.md`.

---

## Adding a new bot — complete checklist

> After completing every step, run `bash scripts/devops-validate.sh`.
> All checks must pass before the work is done.

1. **Create** `cmd/<newbot>/main.go` calling `bot.Run`.
2. **Create** `cmd/<newbot>/CLAUDE.md` with Goals, Features, Dependencies, Edge Cases.
3. **Create** `wiki/bots/<NewBot>.md` documenting the bot.

4. **`docker-compose.yml`** — add a service block:
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

5. **`docker/docker-compose.yml`** — add a service block:
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

6. **`.github/workflows/ci.yml`** — add `<newbot>` to the build matrix:
   ```yaml
   matrix:
     bot: [bluebot, bunkbot, covabot, djcova, ratbot, <newbot>]
   ```

7. **`.github/workflows/main.yml`** — add two entries:
   - In the `paths-filter` block:
     ```yaml
     <newbot>:
       - 'cmd/<newbot>/**'
     ```
   - In the `workflow_dispatch` matrix:
     ```bash
     echo 'bots=["bluebot","bunkbot","covabot","djcova","ratbot","<newbot>"]' >> $GITHUB_OUTPUT
     ```

8. **`scripts/deployment/health-check.sh`** — add `"<newbot>"` to `EXPECTED_SERVICES`:
   ```bash
   EXPECTED_SERVICES=(bluebot bunkbot covabot djcova ratbot <newbot>)
   ```

9. **`AGENTS.md`** — update the bot list in the Architecture and Bots sections.

10. **Run validation**:
    ```bash
    bash scripts/devops-validate.sh
    ```
    Fix every `FAIL` before committing.

11. **Update `wiki/Home.md`** and add a `wiki/bots/<NewBot>.md` page.
12. **Add an entry to `wiki/Changelog.md`**.

---

## Branch protection — `main`

| Rule | Setting |
|---|---|
| Required status checks | `Validate DevOps Consistency`, `Lint`, `Test` |
| Branches must be up to date | Yes (strict mode) |
| Required PR approvals | 1 |
| Dismiss stale reviews on new commits | Yes |
| Force pushes | Blocked |
| Branch deletion | Blocked |

- **Never push directly to `main`.** All changes go through a PR.
- The PR branch must be up to date with `main` before merging.
- If the PR touches any DevOps files, `Validate DevOps Consistency` must pass.

---

## Wiki maintenance

`AGENTS.md` is the single source of truth for rules and architecture. Update it
whenever bots, CI/CD, shared packages, or branch protection rules change.

`CLAUDE.md` imports this file and adds Claude Code-specific notes.
`.github/copilot-instructions.md` points GitHub Copilot here.

Keep all three consistent — `AGENTS.md` is the canonical source.
