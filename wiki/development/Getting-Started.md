# Getting Started

## Prerequisites

- Go 1.21+
- Docker + Docker Compose
- A Discord bot token (set `DISCORD_TOKEN` or `STARBUNK_TOKEN`)

## Run a single bot locally

```bash
DISCORD_TOKEN=<token> go run ./cmd/bunkbot
```

## Run all bots via Docker (local dev)

```bash
docker compose -f docker/docker-compose.yml up -d --build
```

## Build all binaries

```bash
for bot in bluebot bunkbot covabot djcova ratbot; do
  go build -o bin/$bot ./cmd/$bot
done
```

## Run tests

```bash
go test ./...
```

## Validate DevOps consistency

Run this after any bot or CI/CD change:

```bash
bash scripts/devops-validate.sh
```

## See Also

- [[../infrastructure/Configuration|Configuration]] — environment variables
- [[Testing|Testing]] — test guide
- [[CI-CD|CI/CD]] — pipeline overview
