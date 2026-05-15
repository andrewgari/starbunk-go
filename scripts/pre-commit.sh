#!/usr/bin/env bash
# Pre-commit hook: golangci-lint + go build
# Matches the CI linter exactly (golangci-lint v2.12.2, same .golangci.yml).
#
# Install:
#   bash scripts/install-hooks.sh
#
# Requires golangci-lint v2.12.2:
#   mise use golangci-lint@2.12.2    (adds to mise.toml — already done)

set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "$REPO_ROOT"

# --- golangci-lint (covers gofmt, goimports, go vet, and all configured linters) ---
if ! command -v golangci-lint &>/dev/null; then
    echo "⚠️  golangci-lint not found — run: mise install"
    echo "   Falling back to go vet + gofmt only."
    STAGED_GO=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)
    if [ -n "$STAGED_GO" ]; then
        UNFORMATTED=$(echo "$STAGED_GO" | xargs gofmt -l)
        if [ -n "$UNFORMATTED" ]; then
            echo "❌ gofmt: files need formatting:"
            echo "$UNFORMATTED" | sed 's/^/  /'
            exit 1
        fi
    fi
    go vet ./... || exit 1
else
    echo "→ golangci-lint run ./..."
    golangci-lint run ./... || exit 1
fi

# --- go build ---
echo "→ go build ./..."
go build ./... || exit 1

echo "✅ pre-commit checks passed"
