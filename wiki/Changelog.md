# Changelog

Running log of all significant work done on starbunk-go.
Add an entry under today's date for every PR or significant change.

---

## 2026-05-07

- Clarified RatBot purpose across all documentation: RatBot is a rat-themed Secret Santa bot for the guild's "Ratmas" event, not a generic trigger/response bot. Updated `AGENTS.md`, `cmd/ratbot/CLAUDE.md`, and `wiki/bots/RatBot.md`.
- Configured branch protection for `main` (required checks: Validate DevOps Consistency, Lint, Test; 1 PR approval; strict up-to-date; no force push/delete).
- Created `AGENTS.md` as canonical cross-tool agent guide.
- Created `CLAUDE.md` importing `AGENTS.md`; added proactive skill use instructions.
- Created `.github/copilot-instructions.md` for GitHub Copilot.
- Created `wiki/` structure with Home, Changelog, per-bot pages, infrastructure, and development pages.
- Created per-bot `cmd/<bot>/CLAUDE.md` files with goals, features, and edge cases.
