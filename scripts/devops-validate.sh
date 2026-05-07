#!/bin/bash
# Validates that all DevOps files are consistent with the bots defined in cmd/.
#
# Run this script any time you add, remove, or rename a bot, or after editing
# any CI/CD or Docker file. It is also executed as a CI check on every PR.
#
# Usage: bash scripts/devops-validate.sh

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

ERRORS=0

fail() { echo "  FAIL  $1"; ERRORS=$((ERRORS + 1)); }
ok()   { echo "  ok    $1"; }

# ── Discover bots from cmd/ ───────────────────────────────────────────────────
BOTS=()
for dir in cmd/*/; do
  bot=$(basename "$dir")
  if [ -f "cmd/${bot}/main.go" ]; then
    BOTS+=("$bot")
  fi
done

if [ ${#BOTS[@]} -eq 0 ]; then
  echo "ERROR: No bots found under cmd/. Is this the repo root?"
  exit 1
fi

echo "Bots discovered in cmd/: ${BOTS[*]}"
echo ""

# ── Check each file for every bot ────────────────────────────────────────────
for bot in "${BOTS[@]}"; do
  echo "[$bot]"

  # 1. docker-compose.yml (root — production GHCR images)
  if grep -q "starbunk-go-${bot}" docker-compose.yml 2>/dev/null; then
    ok "docker-compose.yml: image starbunk-go-${bot}"
  else
    fail "docker-compose.yml: missing service / image for '${bot}'"
  fi

  # 2. docker/docker-compose.yml (local dev — build from source)
  if grep -q "BOT_NAME: ${bot}" docker/docker-compose.yml 2>/dev/null; then
    ok "docker/docker-compose.yml: BOT_NAME=${bot}"
  else
    fail "docker/docker-compose.yml: missing BOT_NAME: ${bot}"
  fi

  # 3. .github/workflows/ci.yml — build matrix
  if grep -q "${bot}" .github/workflows/ci.yml 2>/dev/null; then
    ok ".github/workflows/ci.yml: build matrix includes ${bot}"
  else
    fail ".github/workflows/ci.yml: missing '${bot}' in build matrix"
  fi

  # 4. .github/workflows/main.yml — paths-filter
  if grep -q "cmd/${bot}" .github/workflows/main.yml 2>/dev/null; then
    ok ".github/workflows/main.yml: paths-filter includes cmd/${bot}"
  else
    fail ".github/workflows/main.yml: missing 'cmd/${bot}' in paths-filter"
  fi

  # 5. scripts/deployment/health-check.sh — EXPECTED_SERVICES
  if grep -qw "${bot}" scripts/deployment/health-check.sh 2>/dev/null; then
    ok "scripts/deployment/health-check.sh: includes ${bot}"
  else
    fail "scripts/deployment/health-check.sh: missing '${bot}' in EXPECTED_SERVICES"
  fi

  echo ""
done

# ── Reverse check: warn about services in compose not backed by a cmd/ dir ───
echo "[reverse check]"
while IFS= read -r svc; do
  # Strip the "starbunk-go-" prefix if present to get the bot name.
  # Services named starbunk-go-<bot> → bot name is after the last dash.
  bot="${svc#starbunk-go-}"
  if [ ! -d "cmd/${bot}" ]; then
    fail "docker-compose.yml: service '${svc}' has no matching cmd/${bot}/"
  else
    ok "docker-compose.yml: service '${svc}' backed by cmd/${bot}/"
  fi
done < <(grep -E '^  [a-z]' docker-compose.yml | grep -v '#' | sed 's/://g' | sed 's/^ *//' || true)

echo ""

# ── Result ────────────────────────────────────────────────────────────────────
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ $ERRORS -gt 0 ]; then
  echo "FAILED — $ERRORS consistency error(s). Fix the files listed above before"
  echo "         committing. See CLAUDE.md § 'DevOps File Maintenance' for the"
  echo "         full list of files that must stay in sync."
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  exit 1
else
  echo "PASSED — all DevOps files are consistent."
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
fi
