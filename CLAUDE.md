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
| User asks about Claude Code / Anthropic API | `claude-code-guide` agent |
| Setting up hooks or automated behaviors | `/update-config` |

@AGENTS.md
