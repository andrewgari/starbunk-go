# CI/CD

## Workflows

### `ci.yml` — Pull Request Checks

Triggered on all PRs to `main`. Jobs:

1. **Validate DevOps Consistency** — runs `scripts/devops-validate.sh` to check if all bots are registered in required config files; skipped if no DevOps files or bot code changed.
2. **Lint** — runs `golangci-lint` on the repository; skipped if no Go files or Go config files changed.
3. **Test** — matrix build; runs `go vet` and `go test` for each affected bot, dynamically running only on packages in the bot's build/runpath; skipped if no bot code changed.
4. **Build** — matrix build; builds each affected bot binary to verify compilation; skipped if no bot code changed.
5. **Docker Test** — matrix build; builds each affected Docker image to verify the Dockerfile; skipped if no bot code changed.

All five jobs are required to pass before a PR can merge (unless they are skipped due to change-detection logic, which counts as passed/success).

#### Selective Validation (Change Detection)

The `ci.yml` workflow optimizes execution times by gating jobs and selectively compiling/testing only the packages affected by a PR's changes:
- **No Go / DevOps changes (e.g. pure documentation/wiki PRs)**: All CI check jobs (DevOps validation, Lint, Test, Build, Docker Test) are skipped.
- **Global / Core Changes**: Modifying `go.mod`, `go.sum`, `docker/Dockerfile`, or core packages (`internal/bot`, `internal/discord`, `internal/middleware`, `internal/other`) triggers builds, vet, tests, and docker checks for **all** bots.
- **Specific Shared Libraries**: Changes under `internal/replybot` only trigger builds/tests/checks for `bluebot`. Changes under `internal/llm` or `internal/memory` only trigger builds/tests/checks for `covabot`.
- **Bot-Specific Code**: Modifying code in a specific bot's directory (e.g. `cmd/bunkbot/`) only triggers builds/tests/checks for that bot.


### `main.yml` — Merge to Main (auto-release)

Triggered on every push to `main`. This is the only workflow that creates releases
and deploys to Tower — **every merge automatically ships**. Jobs:

1. **Validate DevOps Consistency**
2. **Lint**
3. **Test**
4. **Determine Version** — reads the last `v*` git tag and the merge commit title
   to compute the next semver (major/minor/patch via conventional commits).
5. **Docker Publish** — builds all five bots in parallel; pushes `:vX.Y.Z`,
   `:latest`, and `:sha-<short-sha>` to GHCR.
6. **Create Release** — creates a `vX.Y.Z` git tag and GitHub Release, which
   triggers `deploy.yml` automatically.

### `deploy.yml` — Deploy to Tower

Triggered automatically when a GitHub Release is published (i.e., after `main.yml`
completes). Tower deploys `:vX.Y.Z` (the specific version that was just released).
See [[../infrastructure/Deployment|Deployment]].

---

## Version bump rules

The `version` job reads the merge commit title (conventional commits):

| Commit title | Bump |
|---|---|
| `feat!:` or body contains `BREAKING CHANGE` | major |
| `feat:` | minor |
| `fix:`, `chore:`, `refactor:`, anything else | patch |

---

## Definition of Done

A task is **not complete** until:

1. All CI checks pass.
2. The PR has at least one approval and all checks are green.
3. `scripts/devops-validate.sh` exits cleanly (if any DevOps file was touched).
4. The relevant `wiki/` page(s) have been updated.
5. An entry has been added to `wiki/Changelog.md`.

## See Also

- `.github/workflows/`
- [[../infrastructure/Deployment|Deployment]]
- [[../Versioning|Versioning]]
