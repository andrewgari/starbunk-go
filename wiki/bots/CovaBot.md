# CovaBot

> AI personality emulator with LLM-driven responses.

## Goals & Purpose

CovaBot impersonates a specific user's tone and personality in Discord, using an
LLM to generate contextually-aware replies. It is a port of the starbunk-js
CovaBot, which uses Ollama (primary), Gemini, and OpenAI as fallbacks.

## Major Features

- Personality-driven LLM response generation.
- Conversational context modelling.
- Multi-provider LLM support (Ollama → Gemini → OpenAI fallback chain).

## Dependencies & Architecture

- **Entry point:** `cmd/covabot/main.go`
- **Framework:** `internal/bot.Run` + `internal/discord.MessagingService`
- **LLM:** `CLOUD_LLM_PROVIDER` / `LOCAL_LLM_PROVIDER` env vars (not yet wired in Go port)
- API calls must be fully async and timeout-resistant.

## Edge Cases

- All three LLM providers failing simultaneously.
- Rate limits and hallucination management.
- Infinite loops when interacting with other bots (bot must ignore other bot authors).
- Parsing extremely long conversation threads efficiently.

## See Also

- `cmd/covabot/CLAUDE.md`
- [[../infrastructure/Configuration|Configuration]] for LLM env vars
