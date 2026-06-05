---
name: build
description: Verify that all packages compile successfully
---

# Build

Compile all Go packages to catch import or syntax errors fast.

## Instructions

1. Run the build command for all packages:
   ```bash
   go build ./...
   ```

2. If you need to build specific bot binaries to verify them individually:
   ```bash
   for bot in bluebot bunkbot covabot djcova ratbot; do
       go build ./cmd/$bot/ || echo "BROKEN: $bot"
   done
   ```

3. Report the outcome.
