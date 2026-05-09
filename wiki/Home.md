# Starbunk-Go Wiki

> **Second Brain for the Starbunk-Go Discord Bot System**
> Last Updated: 2026-05-07

---

## What Is Starbunk-Go?

Starbunk-Go is a Go monorepo containing 5 independent Discord bots, each with
its own binary and Docker container. It is a port of
[starbunk-js](https://github.com/andrewgari/starbunk-js) to Go, sharing the
same bot personalities but using a simpler, single-binary-per-bot architecture
with no shared database dependencies (yet).

---

## Navigation

### Bots
- [[bots/BlueBot|BlueBot]] — Pattern-matching bot for "blue" / Blue Mage references
- [[bots/BunkBot|BunkBot]] — Administrative backbone and general reply bot
- [[bots/CovaBot|CovaBot]] — AI personality emulator (LLM-driven)
- [[bots/DJCova|DJCova]] — Voice channel music streaming (YouTube)
- [[bots/RatBot|RatBot]] — Rat-trigger response bot

### Infrastructure & Deployment
- [[infrastructure/Architecture|Architecture]] — Monorepo layout, bot pattern, shared libraries
- [[infrastructure/Deployment|Deployment]] — CI/CD pipeline, Docker images, Tower server
- [[infrastructure/Configuration|Configuration]] — Environment variables and Docker Compose files

### Development
- [[development/Getting-Started|Getting Started]] — Local dev setup
- [[development/Testing|Testing]] — Ginkgo/Gomega test guide
- [[development/CI-CD|CI/CD]] — GitHub Actions workflows
- [[development/MessageFiltering|Message Filtering]] — Composable message evaluation abstraction (planned)

### AI Agents
- [[agents/Agents|Custom Agents]] — Claude Code subagents: go-craftsman, architect, pm, devops

### History
- [[Changelog|Changelog]] — Running log of all work done on this project

---

## Quick Reference

| Bot | Image | Binary |
|-----|-------|--------|
| BlueBot | `ghcr.io/andrewgari/starbunk-go-bluebot` | `cmd/bluebot` |
| BunkBot | `ghcr.io/andrewgari/starbunk-go-bunkbot` | `cmd/bunkbot` |
| CovaBot | `ghcr.io/andrewgari/starbunk-go-covabot` | `cmd/covabot` |
| DJCova  | `ghcr.io/andrewgari/starbunk-go-djcova`  | `cmd/djcova`  |
| RatBot  | `ghcr.io/andrewgari/starbunk-go-ratbot`  | `cmd/ratbot`  |

### Key Commands

```bash
go test ./...                              # run all tests
bash scripts/devops-validate.sh            # validate DevOps file consistency
docker compose -f docker/docker-compose.yml up -d --build   # local dev
```

---

## Agent Instructions

> These instructions apply to **all AI agents** working in this repository.
> They are also codified in `AGENTS.md` at the repo root.

1. **Before starting any task** — read the relevant wiki page(s) for the area you will touch.
2. **After completing any task** — update the relevant wiki page(s) with any changes to architecture, behavior, config, or patterns.
3. **For every significant change or PR** — add an entry to [[Changelog]] under today's date.
4. If a wiki page does not exist for the area you are working in, create it.
5. Use Obsidian-style `[[Page]]` or `[[folder/Page|Display Name]]` links between pages.
