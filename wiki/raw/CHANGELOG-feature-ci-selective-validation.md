## [Unreleased] — feature/ci-selective-validation

### Changed
- Optimized Pull Request CI pipeline (`ci.yml`) to selectively run build and smoke test checks only for bots affected by changes.
- Mapped internal package dependencies to their specific importing bots (e.g. `internal/replybot` changes trigger `bluebot`; `internal/llm`/`memory` trigger `covabot`).
- Documented the change-detection strategy in `wiki/development/CI-CD.md`.
