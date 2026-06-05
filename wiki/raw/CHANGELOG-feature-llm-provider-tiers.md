## [Unreleased] — feature/llm-provider-tiers

### Added
- Added `Embed` to the `llm.Service` interface for generating vector embeddings.
- Added `ClientConfig` and global `Config` struct to `internal/llm/` to support configuring multiple providers.
- Added `Registry` to `internal/llm/` which acts as a factory providing `High()`, `Medium()`, and `Low()` tier LLM clients based on configuration.

### Changed
- Refactored `llm` package to no longer use a single environment variable prefix (`CLOUD_LLM_` or `LOCAL_LLM_`) and instead load all providers (OpenAI, Anthropic, Google, Ollama) if credentials exist in `.env`.
- `new[Provider]Client` constructors now accept `ClientConfig` instead of `Config`.
