## [Unreleased] — feat-agent-workflows

### Added
- Ported and adapted AI developer skills and workflows from `starbunk-js` to `starbunk-go`.
- Created `.gemini/skills/` directory housing Go-specific development skills: `check`, `lint`, `test`, `build`, `health-check`, `git-workflow`, `review`.
- Created custom AI Agents and commands for both Gemini and Claude (`task-runner`, `address-pr-comments`, `debugger`, and `ci-diagnose`) to handle complex interactive workflows and automated troubleshooting.
- Created `GEMINI.md` to specify proactive skill usage and agentic rules for Gemini.
- Updated `.claude/settings.json` to add a `postToolUse` hook running `go build ./...` after edits/writes.
- Added documentation for the new dual-agent (Claude Code + Gemini) ecosystem in `wiki/agents/Agents.md` and `AGENTS.md`.
