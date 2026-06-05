---
name: check
description: Run the full CI test suite locally (go test, go vet, golangci-lint, and devops-validate)
---

# Full Check

Run all local CI checks to ensure code is ready for PR.

## Instructions

Run the following checks sequentially:

1. **DevOps Validation**:
   ```bash
   bash scripts/devops-validate.sh
   ```
2. **Go Vet**:
   ```bash
   go vet ./...
   ```
3. **Linting**:
   ```bash
   golangci-lint run
   ```
4. **Testing**:
   ```bash
   go test ./...
   ```

Report the results of each step clearly. If any step fails, analyze the output and suggest fixes.
