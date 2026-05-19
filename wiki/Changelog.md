# Changelog

Running log of all significant work done on starbunk-go.
Add an entry under today's date for every PR or significant change.

---

## 2026-05-14 — Add self-correction protocol to AGENTS.md

### Added
- Self-Correction Protocol section in `AGENTS.md` with six recovery steps:
  0. pre-edit sync/read/blast-radius check
  1. post-change verification loop (`go build`, `go vet`, `go test`, per-bot build, `golangci-lint`)
  2. error-reading discipline (no blind retries, no `//nolint` suppression)
  3. DevOps drift check trigger rule
  4. wiki/code consistency cross-reference table
  5. CI failure recovery sequence
  6. uncertainty protocol (build → wiki → git log → ask)

---

## 2026-05-14 — Port agentic guidance from starbunk-js

### Added
- CHANGELOG staging workflow (`wiki/raw/CHANGELOG-<branch>.md` pattern) in `AGENTS.md`.
- Development Constraints section in `AGENTS.md`: no secret commits, container
  isolation rules, self-message guard, non-blocking handler requirement.
- "What counts as meaningful" definition for wiki update rule in `AGENTS.md`.
- Pre-done checklist: `go test ./...` locally required before declaring tasks done.
- BlueBot rate-limit specifics (5-min standard / 24-hour rare reply windows).
- CovaBot 4-stage decision pipeline and social battery documented in `cmd/covabot/CLAUDE.md`.
- LLM provider priority order (Ollama → Anthropic → Gemini → OpenAI) added to CovaBot docs.

---

## 2026-05-19 (3)

- **Fully automated release pipeline**: every merge to `main` now auto-versions,
  builds, and deploys — no manual `git tag` needed.
  - `main.yml` bumps semver from conventional commit title (feat→minor,
    feat!→major, everything else→patch), builds all 5 bots, pushes
    `:vX.Y.Z` + `:latest` + `:sha-*`, creates git tag + GitHub Release.
  - `release.yml` deleted — it is no longer needed.
  - `deploy.yml` now pins Tower to the specific version tag (`:vX.Y.Z`)
    rather than always pulling `:latest`.

## 2026-05-19 (2)

- **Compose/deploy audit guidance added to AGENTS.md**: added "Image tag chain
  audit" checklist to the DevOps maintenance section, requiring agents to
  verify workflow image names, tag variables, and pre-release behaviour when
  editing CI/CD workflows.
- **Fixed wiki inaccuracy** in `wiki/Versioning.md`: pre-release deploys do
  _not_ pull the RC image tag — Tower continues running `:latest`. Corrected
  the description and documented the manual override process.

---

## 2026-05-15 (3)

- **Tag-based release system**: replaced automatic semver-from-PR-labels with
  an explicit `git tag v1.3.0 && git push origin v1.3.0` release flow.
  - New `release.yml` workflow: triggers on `v*` tag push, validates the tag
    is on `main`, runs lint+test, builds all 5 bot images tagged `:vX.Y.Z` +
    `:latest`, creates GitHub Release (which triggers Tower deploy).
  - Simplified `main.yml`: removed `semver_tag`, `publish_release`, and
    `versioned_bots` logic. Now only publishes `:main`/`:sha-*` images and
    creates `build-YYYYMMDD-sha` breadcrumb tags on merge.
  - `:latest` is now exclusively owned by `release.yml` — Tower always runs a
    named, intentional release.

## 2026-05-15 (2)

- Migrated `internal/bluebot/` into `cmd/bluebot/` — `BlueStrategy` and its
  Ginkgo test suite now live alongside `main.go` as `package main`.
  The `internal/bluebot` package is deleted; no other packages are affected.

## 2026-05-15

- Added two-tier condition system to `internal/replybot`:
  - **Tier 1** (unchanged): bot-level `MessageAuditor` in `bot.Run()` — hard gates like `NotSelf`.
  - **Tier 2** (new): strategy-level conditions via the optional `ConditionedStrategy` interface
    and the `WithCondition(cond, strategy)` compose helper in `replybot/condition.go`.
  - `Bot.Handle()` now accepts `*discordgo.Session` so strategy conditions (e.g. `AuthorHasRole`)
    can inspect guild state.
  - Enables BunkBot to host mixed strategies: `WithCondition(middleware.IsBot, botBotStrategy)`
    alongside `WithCondition(middleware.NotBot, humanOnlyStrategy)`.
  - Added Ginkgo test suite `internal/replybot/bot_test.go` covering all dispatch cases.

## 2026-05-14

- Fixed critical bug in `.github/workflows/ci.yml`: `docker_test` job was
  pushing `:latest` and `sha-{hash}` images to GHCR on every PR, overwriting
  the production tag. CI now builds and smoke-tests Docker images locally only;
  GHCR publishing happens exclusively in `main.yml` after merge.
- Removed `packages: write` permission from `ci.yml` (PRs never need to write
  to GHCR).
- Added `.golangci.yml` with an explicit linter set (`govet`, `errcheck`,
  `staticcheck`, `ineffassign`, `unused`, `gofmt`, `goimports`, `misspell`,
  `whitespace`, `gosec`, `noctx`) so linter behaviour is pinned and
  reproducible across golangci-lint versions.

## 2026-05-13 (2)

- Deleted `internal/config.go` and `internal/config_test.go` — dead scaffold
  from project init (`package pkg`, `func HelloShared()`).
- Added `/bluebot`, `/bunkbot`, `/covabot`, `/djcova`, `/ratbot` to
  `.gitignore` so root-level built binaries are never accidentally staged.
- Tightened `scripts/devops-validate.sh` ci.yml check: changed
  `grep "${bot}"` to `grep "cmd/${bot}/"` so a bot name in a comment no
  longer satisfies the check.
- Migrated all logging from `log`/`fmt` to `log/slog` (stdlib, Go 1.21).
  Every log call now emits structured fields: `bot=`, `strategy=`,
  `channel=`, `err=`. No new dependencies.

## 2026-05-13

- Implemented BlueBot strategy engine in `internal/bluebot/` and shared
  dispatcher in `internal/replybot/`.
- New `Strategy` interface (`Name`, `ShouldTrigger(ctx, msg)`, `Response(ctx, msg)`)
  in `internal/replybot/` is the extensibility seam for all reply-style bots.
  `ctx` is threaded through every call so LLM-backed strategies can respect
  cancellation deadlines when that time comes.
- `Bot` dispatcher in `internal/replybot/`: ordered strategies, first match
  wins. `Handle(ctx, m)` — caller supplies the context, no hardcoded
  `context.Background()` buried inside.
- `BlueStrategy` in `internal/bluebot/`: regex covering `blue`, `blu+`,
  `bloo+`, `blew`, `bleu`, `azul`, `blau`, `bluebot` (case-insensitive,
  word-bounded). Response: `"Did somebody say Blu?"`. Gains compile-time
  interface assertion `var _ replybot.Strategy = BlueStrategy{}`.
- `cmd/bluebot/main.go`: stub replaced with real engine; `Bot` constructed
  once via `sync.Once`; calls `blueBot.Handle(context.Background(), m)`.
- 25 Ginkgo specs: matches, case variants, compound-word exclusions
  (`bluetooth`, `blueprint`, `blueberry`), catchphrase.
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
