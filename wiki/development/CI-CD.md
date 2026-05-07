# CI/CD

## Workflows

### `ci.yml` — Pull Request Checks

Triggered on all PRs to `main`. Jobs:

1. **Validate DevOps Consistency** — runs `scripts/devops-validate.sh`; fails fast if any bot is not registered in all required files.
2. **Lint** — runs `golangci-lint`.
3. **Test** — runs `go vet` and `go test ./...`.
4. **Build** — matrix build; builds each bot binary to verify compilation.
5. **Docker Test** — builds each Docker image to verify the Dockerfile.

All five jobs are required to pass before a PR can merge.

### `main.yml` — Merge to Main

Triggered on push to `main`. Jobs:

1. **Validate DevOps Consistency**
2. **Test**
3. **Build + Push Docker Images** — builds and pushes each bot image to GHCR; tags `:latest`.
4. **Tag Release** — creates a git build tag.

### `deploy.yml` — Deploy to Tower

Triggered automatically after `main.yml` succeeds. See [[../infrastructure/Deployment|Deployment]].

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
