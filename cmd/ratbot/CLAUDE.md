# RatBot — Development Instructions

> See also: `wiki/bots/RatBot.md`

## Goals & Purpose

Monitors guild messages for rat-related triggers and replies accordingly.
Specific trigger patterns and responses are TBD; the bot follows the same
pattern-matching architecture as BlueBot.

## Major Features

- Pattern matching for rat-related keywords or phrases.
- Configurable response content.

## Dependencies & Architecture

- `internal/bot.Run` — event loop and session management.
- `internal/discord.MessagingService` — sends replies.
- No external services required.

## Edge Cases

- Avoid false positives on messages that contain "rat" as a substring of unrelated words.
- Rate-limit replies to prevent channel spam.
