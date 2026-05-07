# Changelog

Running log of all significant work done on starbunk-go.
Add an entry under today's date for every PR or significant change.

---

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
