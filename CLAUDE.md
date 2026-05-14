# CLAUDE.md

Agent guide for Claude Code. All rules, architecture notes, and the DevOps
maintenance checklist live in [AGENTS.md](AGENTS.md) and are imported below.

## Claude Code: Proactive Skill Use

Use available skills **without being told**. When the situation matches, invoke
the skill immediately — don't describe what you'd do, just do it.

| Situation | Skill |
|---|---|
| Any coding, fixing, or refactoring task | `/task` |
| Code has been written or changed | `/simplify` — review for quality and reuse |
| Deploying or updating containers on Tower | `/deploy` |
| PR is open — review comments to address | run `go test ./...` locally, then address each comment |
| User asks about Claude Code / Anthropic API | `claude-code-guide` agent |
| Setting up hooks or automated behaviors | `/update-config` |

Before declaring any task done, run `go test ./...` locally. If tests fail,
fixing them is part of the task.

@AGENTS.md
