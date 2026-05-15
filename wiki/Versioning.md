# Versioning

Versions in starbunk-go are **explicit git tags** pushed by the developer.
The pipeline reacts to those tags — it does not invent versions automatically.

---

## Release flow

```bash
git checkout main && git pull
git tag v1.3.0 -m "Release v1.3.0"
git push origin v1.3.0
```

Pushing the tag triggers `release.yml`, which:

1. Validates the tag points to a commit on `main`.
2. Runs lint and tests.
3. Builds all five bot images, pushes `:v1.3.0` and `:latest` to GHCR.
4. Creates a GitHub Release (which triggers Tower deployment via `deploy.yml`).

---

## Tag format

Repo-level semver tags:

```
v<MAJOR>.<MINOR>.<PATCH>
```

Examples: `v1.2.3`, `v0.3.0`, `v1.0.0-rc.1`

Per-bot tags (legacy, from pre-2026-05 pipeline) also exist in the tag history
(`bluebot/v0.2.4` etc.) but are no longer created by the new release workflow.

---

## Pre-release tags

Tags with a hyphen suffix are treated as pre-releases:

```bash
git tag v1.3.0-rc.1 -m "Release candidate"
git push origin v1.3.0-rc.1
```

- Images are tagged `:v1.3.0-rc.1` but **not** `:latest`.
- The GitHub Release is marked `--prerelease`.
- `deploy.yml` still fires (Tower pulls `:v1.3.0-rc.1`); add an environment
  protection rule in `deploy.yml` if you want manual approval for RCs.

---

## Docker image tags

| Tag | Set by | Meaning |
|---|---|---|
| `:latest` | `release.yml` only | Most recent stable release |
| `:v<MAJOR>.<MINOR>.<PATCH>` | `release.yml` | Specific versioned release |
| `:main` | `main.yml` | Current HEAD of `main` (not deployed) |
| `:sha-<content-hash>` | `main.yml` / `release.yml` | Content-addressed; unchanged if source didn't change |
| `:sha-<short-sha>` | `main.yml` | Commit SHA |

Tower's `docker-compose.yml` uses `${IMAGE_TAG:-latest}`, so it runs the
last explicitly released version by default.

---

## Finding the current version

```bash
# Latest release tag
git tag -l 'v[0-9]*.[0-9]*.[0-9]*' | grep -vE -- '-' | sort -V | tail -1

# All release tags
git tag -l 'v[0-9]*.[0-9]*.[0-9]*' | sort -V
```

---

## See Also

- [[development/CI-CD|CI/CD]] — workflow details
- [[infrastructure/Deployment|Deployment]] — how Tower pulls the right image tag
