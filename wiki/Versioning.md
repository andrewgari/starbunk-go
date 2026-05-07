# Versioning

Each bot in starbunk-go has its own independent **semantic version** (semver).
Versions are stored as **git tags** on the merge commit and as **Docker image tags** in GHCR.

---

## Tag format

```
<bot>/v<MAJOR>.<MINOR>.<PATCH>
```

Examples: `bluebot/v1.2.3`, `bunkbot/v0.3.0`, `ratbot/v0.0.1`

Docker images are tagged both `v<MAJOR>.<MINOR>.<PATCH>` and `:latest`.

---

## How versioning works

When a PR is merged to `main`, the `CD / Main Merge` workflow:

1. Detects which bots have code changes (same paths-filter as the build step).
2. Reads the **PR labels** to determine the bump type.
3. Creates a `<bot>/vX.Y.Z` git tag on the merge commit for every affected bot.
4. Re-tags the published Docker images with the semver version.

CI-only PRs (e.g., editing only `.github/workflows/`) trigger safety-net
rebuilds but **do not bump any bot version** — only actual code changes count.

---

## Bump type resolution

The bump type is determined in this order of precedence:

| Priority | Trigger | Bump |
|---|---|---|
| 1 | PR label: `bump:major` or `breaking` | major (`X+1.0.0`) |
| 2 | PR label: `bump:minor`, `feat`, `feature`, or `enhancement` | minor (`X.Y+1.0`) |
| 3 | PR label: `bump:patch`, `fix`, `bug`, or `chore` | patch (`X.Y.Z+1`) |
| 4 | PR title matches `feat(...)!:` or contains `BREAKING CHANGE` | major |
| 5 | PR title matches `feat(...):` | minor |
| 6 | *(default)* | patch |

---

## PR labels

Apply one of these labels to a PR before merging:

| Label | Effect |
|---|---|
| `bump:major` | Major version bump — breaking change |
| `bump:minor` | Minor version bump — new feature |
| `bump:patch` | Patch version bump — bug fix or chore |
| `breaking` | Alias for `bump:major` |

If no label is applied, the workflow falls back to conventional commit title
inference (step 4–5 above), then defaults to a **patch** bump.

### Seeding labels

Run these once to create the labels in the GitHub repo:

```bash
gh label create "bump:major" --color "d73a4a" --description "Breaking change — bumps major version (X.0.0)"
gh label create "bump:minor" --color "0075ca" --description "New feature — bumps minor version (0.X.0)"
gh label create "bump:patch" --color "0e8a16" --description "Bug fix or chore — bumps patch version (0.0.X)"
gh label create "breaking"   --color "b60205" --description "Contains a breaking change (alias for bump:major)"
```

---

## Manual trigger

The `CD / Main Merge` workflow can be dispatched manually from the GitHub UI
or via the CLI. A `bump_type` input lets you choose the bump level:

```bash
gh workflow run main.yml --field bump_type=minor
```

Valid values: `patch` (default), `minor`, `major`.

---

## Shared library changes

When `internal/**`, `go.mod`, `go.sum`, or `docker/Dockerfile` change, **all
bots** are rebuilt and version-bumped together, using the same bump type
determined from the PR labels.

---

## Docker image tags

For each affected bot after a merge to `main`:

| Tag | Meaning |
|---|---|
| `:latest` | Most recent build |
| `:main` | Current HEAD of `main` |
| `:v<MAJOR>.<MINOR>.<PATCH>` | Specific semver release |
| `:sha-<hash>` | Content-addressed (unchanged if source didn't change) |
| `:sha-<short-sha>` | Commit SHA |

---

## Finding current versions

```bash
# List all version tags for a specific bot
git tag -l "bluebot/v*" | sort -V

# Latest version of every bot
for bot in bluebot bunkbot covabot djcova ratbot; do
  latest=$(git tag -l "${bot}/v*" | grep -E "^${bot}/v[0-9]+\.[0-9]+\.[0-9]+$" | sort -V | tail -1)
  echo "${bot}: ${latest:-no version yet}"
done
```
