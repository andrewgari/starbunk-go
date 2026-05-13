# Changelog

Running log of all significant work done on starbunk-go.
Add an entry under today's date for every PR or significant change.

---

## 2026-05-13

- Implemented BlueBot strategy engine in `internal/bluebot/`.
- New `Strategy` interface (`Name`, `ShouldTrigger`, `Response`) is the
  extensibility seam — swap in an LLM or add stateful sub-strategies without
  touching the dispatcher.
- `BlueStrategy`: regex-based trigger covering `blue`, `blu+`, `bloo+`, `blew`,
  `bleu`, `azul`, `blau`, `bluebot` (case-insensitive, word-bounded).
  Response: `"Did somebody say Blu?"`.
- `Bot` dispatcher: holds an ordered slice of strategies, first match wins.
  Constructed once via `sync.Once` so future stateful strategies preserve
  their state across calls.
- `cmd/bluebot/main.go` updated — stub ping handler replaced with real engine.
- 25 Ginkgo specs covering: matches, case-insensitivity, compound-word
  exclusions (`bluetooth`, `blueprint`, `blueberry`), and the catchphrase.
- Updated `wiki/bots/BlueBot.md` with architecture overview and extensibility
  roadmap.

## 2026-05-09

- Added `internal/middleware` package: composable `MessageAuditor` interface for mandatory message evaluation.
- Primitives split by concern across `author.go`, `content.go`, `context.go`, `random.go`.
- Gates: `NotSelf`, `NotBot`, `IsBot`, `AuthorID`, `NotAuthorID`, `AuthorNamed`, `AuthorHasRole`, `HasContent`, `ContentContains`, `ContentMatches`, `HasAttachment`, `GuildOnly`, `DMOnly`, `InChannel`, `OnWeekdays`, `Chance`.
- Combinators: `AllOf`, `AnyOf`, `Not`.
- `bot.Run` signature updated to require a `MessageAuditor`; all `MessageCreate` handlers are automatically wrapped — no handler can be invoked without passing audit.
- All 5 bots (`bluebot`, `bunkbot`, `covabot`, `djcova`, `ratbot`) updated to declare their own auditor in `main.go`; inline self-guards removed.
- 48 Ginkgo specs covering every primitive, combinator, and both complex composition scenarios.
- Updated `wiki/infrastructure/Architecture.md` and `wiki/development/MessageFiltering.md`.

## 2026-05-07 (6)

- `deploy.yml` now triggers on `release: published` instead of `workflow_run`.
- `main.yml` gains a new `publish_release` job that creates a GitHub Release (using the build tag) once Docker images are successfully published and semver tags are created. The release notes list each changed bot and its new version.
- Deploy uses the `:latest` image tag (was `:main`) and passes the release tag name as the version label.
- This ensures deployment happens exactly once per PR merge that produces new images, and only after all upstream jobs succeed.

## 2026-05-07 (5)

- Added 4 custom Claude Code subagents to `.claude/agents/`:
  - `go-craftsman` — Go code writing, idiomatic patterns, naming, aesthetics
  - `architect` — High-level planning, cross-cutting concerns, directing other agents (uses Opus)
  - `pm` — Requirements gathering, clarifying questions, scope alignment
  - `devops` — GitHub, CI/CD, Docker, deployment management
- Added `wiki/agents/Agents.md` documenting the agents and when to use each.

## 2026-05-07 (4)

- Fixed `main.yml`: `semver_tag` and `tag_release` jobs now depend on `validate_devops`, `lint`, and `test` directly instead of `docker_publish`. This ensures git tags are always created when code changes and tests pass, even if Docker image publishing partially fails for one or more bots.

## 2026-05-07 (3)

- Added per-bot semver versioning: each bot (`bluebot`, `bunkbot`, `covabot`, `djcova`, `ratbot`) now gets an independent `<bot>/vX.Y.Z` git tag on every merge to `main`.
- Bump type is driven by PR labels (`bump:major`, `bump:minor`, `bump:patch`, `breaking`) with fallback to conventional commit title inference.
- CI-only PRs (no bot code changed) are skipped — only real code changes trigger version bumps.
- Docker images are also tagged with the semver version (`:v<X.Y.Z>`) in addition to `:latest` / `:main`.
- `workflow_dispatch` gains a `bump_type` input (`patch` / `minor` / `major`) for manual triggers.
- Added `.github/labels.yml` documenting the version labels and seed commands.
- Added `wiki/Versioning.md` documenting the full versioning scheme.

## 2026-05-07 (2)

- Clarified RatBot purpose across all documentation: RatBot is a rat-themed Secret Santa bot for the guild's "Ratmas" event, not a generic trigger/response bot. Updated `AGENTS.md`, `cmd/ratbot/CLAUDE.md`, and `wiki/bots/RatBot.md`.

## 2026-05-07

- Clarified RatBot purpose across all documentation: RatBot is a rat-themed Secret Santa bot for the guild's "Ratmas" event, not a generic trigger/response bot. Updated `AGENTS.md`, `cmd/ratbot/CLAUDE.md`, and `wiki/bots/RatBot.md`.
- Configured branch protection for `main` (required checks: Validate DevOps Consistency, Lint, Test; 1 PR approval; strict up-to-date; no force push/delete).
- Created `AGENTS.md` as canonical cross-tool agent guide.
- Created `CLAUDE.md` importing `AGENTS.md`; added proactive skill use instructions.
- Created `.github/copilot-instructions.md` for GitHub Copilot.
- Created `wiki/` structure with Home, Changelog, per-bot pages, infrastructure, and development pages.
- Created per-bot `cmd/<bot>/CLAUDE.md` files with goals, features, and edge cases.
