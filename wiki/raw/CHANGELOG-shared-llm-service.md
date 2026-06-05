## [Unreleased] — shared-llm-service

### Added
- Created `internal/llm` package to serve as a shared abstraction for bots to communicate with LLMs.
- Defined `GenerateRequest`, `GenerateResponse`, and `ResponseSchema` models (validation enforcement to be implemented).
- Defined `Service` interface for agnostic provider support (OpenAI, Anthropic, Ollama, Google).
- Implemented raw HTTP clients using `net/http` to strictly map JSON payloads to each provider's API.
- Defined `Config` struct and `ConfigFromEnv` for configuring URL, port, provider, and API keys via environment variables (e.g., `CLOUD_LLM_PROVIDER`).
