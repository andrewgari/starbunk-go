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

Triggered on push to `main`. Publishes continuous integration images but does **not** create a release or deploy to Tower. Jobs:

1. **Detect Changed Bots** — determines which bots need rebuilding.
2. **Validate DevOps Consistency**
3. **Lint**
4. **Test**
5. **Docker Publish** — builds and pushes `:main` and `:sha-*` images to GHCR for changed bots.
6. **Tag Release** — creates a `build-YYYYMMDD-sha` git tag as a breadcrumb.

### `release.yml` — Tag Push (Release)

Triggered when a `v*` tag is pushed to the repo (e.g., `git tag v1.3.0 && git push origin v1.3.0`). This is the only workflow that deploys to Tower. Jobs:

1. **Validate Tag** — confirms the tag points to a commit on `main` and that no release for this version already exists.
2. **Lint**
3. **Test**
4. **Docker Release** — builds all five bots, pushes `:v1.3.0` and `:latest` images to GHCR.
5. **Publish Release** — creates a GitHub Release, which triggers `deploy.yml`.

### `deploy.yml` — Deploy to Tower

Triggered automatically when a GitHub Release is published (i.e., after `release.yml` completes). See [[../infrastructure/Deployment|Deployment]].

---

## Release workflow

```bash
# Merge PRs to main normally — nothing deploys automatically.

# When you're ready to ship:
git checkout main && git pull
git tag v1.3.0 -m "Release v1.3.0"
git push origin v1.3.0
# → release.yml runs → images pushed → GitHub Release created → Tower deploys
```

Pre-release tags (e.g., `v1.3.0-rc.1`) are supported — they publish images but
mark the GitHub Release as a pre-release and do not update `:latest`.

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
