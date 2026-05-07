# Configuration

## Environment Variables

| Variable | Purpose |
|---|---|
| `DISCORD_TOKEN` | Token used by `bot.Run` at runtime |
| `STARBUNK_TOKEN` | Fallback token used by all bots in Docker Compose |
| `{BOTNAME}_TOKEN` | Per-bot override (e.g. `BUNKBOT_TOKEN`, `COVABOT_TOKEN`) |
| `CLOUD_LLM_PROVIDER` | Cloud LLM provider name (e.g. `gemini`) — not yet wired |
| `CLOUD_LLM_API_KEY` | API key for cloud LLM |
| `LOCAL_LLM_PROVIDER` | Local LLM provider name (e.g. `ollama`) — not yet wired |
| `LOCAL_LLM_API_KEY` | API key / endpoint for local LLM |

Each Docker Compose service resolves its token as:
```
${BOTNAME_TOKEN:-${STARBUNK_TOKEN}}
```

## Docker Compose Files

| File | Purpose |
|---|---|
| `docker-compose.yml` | **Production** — pulls pre-built GHCR images. Deployed to Tower by `deploy.yml`. Requires `stack.env` on the server. |
| `docker/docker-compose.yml` | **Local dev** — builds from source using `docker/Dockerfile` with `BOT_NAME` build arg. |

## Local Dev Setup

```bash
cp .env.example .env   # if present, or create manually
# set STARBUNK_TOKEN=<your dev token>
docker compose -f docker/docker-compose.yml up -d --build
```

## See Also

- [[../development/Getting-Started|Getting Started]]
