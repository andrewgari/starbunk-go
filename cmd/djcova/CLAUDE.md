# DJCova — Development Instructions

> See also: `wiki/bots/DJCova.md`

## Goals & Purpose

Voice channel music streaming service. Joins Discord voice on demand, streams
YouTube audio, and manages a per-guild playback queue. Ported from starbunk-js
DJCova.

## Major Features

- `/play <youtube-url>` — joins voice channel and streams audio.
- Per-guild queue (add, skip, clear).
- Voice channel state management (join, leave, idle timeout, reconnect).

## Dependencies & Architecture

- `internal/bot.Run` — **must be extended with voice intents** (`IntentsGuildVoiceStates`).
- `internal/discord.MessagingService` — status replies in text channel.
- Audio pipeline: discordgo voice + ffmpeg (not yet wired in Go port).
- CPU-intensive audio processing — must not block the event loop; use goroutines.

## Edge Cases

- Monitor voice connection health; reconnect on disconnect.
- Concurrent `/play` requests: serialize queue writes.
- YouTube geo-restrictions and playback errors — report to text channel, skip to next.
- Clean up ffmpeg processes on bot disconnect or crash (defer + signal handling).
