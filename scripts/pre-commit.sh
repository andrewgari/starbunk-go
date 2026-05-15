#!/usr/bin/env bash
# Pre-commit hook: gofmt, go vet, go build
# Install: cp scripts/pre-commit.sh .git/hooks/pre-commit && chmod +x .git/hooks/pre-commit
#   or run: bash scripts/install-hooks.sh

set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "$REPO_ROOT"

# --- gofmt: check staged .go files only ---
STAGED_GO=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

if [ -n "$STAGED_GO" ]; then
    UNFORMATTED=$(echo "$STAGED_GO" | xargs gofmt -l)
    if [ -n "$UNFORMATTED" ]; then
        echo "❌ gofmt: the following staged files need formatting:"
        echo "$UNFORMATTED" | sed 's/^/  /'
        echo ""
        echo "Run: gofmt -w \$(git diff --cached --name-only --diff-filter=ACM | grep '\\.go\$')"
        exit 1
    fi
fi

# --- go vet ---
echo "→ go vet ./..."
if ! go vet ./...; then
    echo "❌ go vet failed"
    exit 1
fi

# --- go build ---
echo "→ go build ./..."
if ! go build ./...; then
    echo "❌ go build failed"
    exit 1
fi

echo "✅ pre-commit checks passed"
