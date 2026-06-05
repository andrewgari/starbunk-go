---
name: test
description: Run the Go test suite for the starbunk-go monorepo or a specific package
---

# Test Runner

Run tests for the project.

## Arguments
- `$ARGUMENTS` - Optional: specific package path (e.g. `./internal/...`) or leave blank for all.

## Instructions

1. If no argument is provided, run the full test suite:
   ```bash
   go test ./...
   ```

2. If an argument is provided, run tests for that package:
   ```bash
   go test <package-path>
   ```

3. Report the test results clearly, highlighting any failures.

4. If tests fail, analyze the output and suggest fixes if the errors are straightforward.
