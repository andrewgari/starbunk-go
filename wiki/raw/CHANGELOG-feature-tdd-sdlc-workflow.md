## [Unreleased] — feature/tdd-sdlc-workflow

### Added
- Documented a new mandatory Test-Driven Development (TDD) SDLC workflow.
- Established a two-PR sequence constraint where tests are added/improved first in a Test-Only PR before the implementation PR is created.
- Created `wiki/development/TDD.md` to guide agents and human developers on how to locate existing behavior in `starbunk-js` and write Ginkgo tests in Go.

### Changed
- Updated `AGENTS.md`, `CLAUDE.md`, and `.github/copilot-instructions.md` to mandate the TDD and two-PR workflow for all AI agents.
- Linked the TDD guide from `wiki/Home.md` and `wiki/development/Testing.md`.
