# CI/CD

## Workflows

### `ci.yml` — Pull Request Checks

Triggered on all PRs to `main`. Jobs:

1. **Validate DevOps Consistency** — runs `scripts/devops-validate.sh`; fails fast if any bot is not registered in all required files.
2. **Lint** — runs `golangci-lint`.
3. **Test** — runs `go vet` and `go test ./...`.
4. **Build** — matrix build; builds each affected bot binary to verify compilation.
5. **Docker Test** — builds each affected Docker image to verify the Dockerfile.

All five jobs are required to pass before a PR can merge.

#### Selective Validation (Change Detection)

The `ci.yml` workflow optimizes build and smoke test times by selectively compiling and testing only the bots affected by a PR's changes:
- **Global / Core Changes**: Modifying `go.mod`, `go.sum`, `docker/Dockerfile`, or core packages (`internal/bot`, `internal/discord`, `internal/middleware`) triggers builds and tests for **all** bots.
- **Specific Shared Libraries**: Changes under `internal/replybot` only trigger builds/tests for `bluebot`. Changes under `internal/llm` or `internal/memory` only trigger builds/tests for `covabot`.
- **Bot-Specific Code**: Modifying code in a specific bot's directory (e.g. `cmd/bunkbot/`) only triggers builds/tests for that bot.


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
