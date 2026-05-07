# CovaBot — Development Instructions

> See also: `wiki/bots/CovaBot.md`

## Goals & Purpose

AI personality emulator. Generates LLM-driven replies that mimic a specific
user's tone and conversational style. Ported from starbunk-js CovaBot.

## Major Features

- Personality-driven LLM response generation.
- Conversational context modelling (conversation history as prompt context).
- Multi-provider LLM fallback chain: Ollama → Gemini → OpenAI.

## Dependencies & Architecture

- `internal/bot.Run` — event loop and session management.
- `internal/discord.MessagingService` — sends replies.
- LLM provider: configured via `LOCAL_LLM_PROVIDER` / `CLOUD_LLM_PROVIDER` env vars (not yet wired).
- All LLM calls must be fully async and timeout-resistant.

## Edge Cases

- Handle all LLM providers failing simultaneously — log and stay silent rather than crash.
- Ignore messages from other bot authors to prevent infinite reply loops.
- Manage LLM rate limits with backoff; do not block the event loop.
- Avoid hallucination bleed-through — truncate or summarise long context threads.
