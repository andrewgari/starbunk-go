# Testing

## Framework

Tests use **Ginkgo v2 / Gomega** BDD framework.

- Test files use the `_test` package suffix (e.g. `package bot_test`).
- Each package needs a suite bootstrap to integrate with `go test`.

## Suite Bootstrap

Add this once per package (e.g. `internal/bot/suite_test.go`):

```go
package bot_test

import (
    "testing"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestBot(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Bot Suite")
}
```

## Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/bot/...

# Single spec (by description)
go test ./internal/... -run "TestInternal/Config"
```

## Writing Specs

```go
var _ = Describe("MessagingService", func() {
    Context("when sending a message", func() {
        It("returns no error", func() {
            Expect(err).NotTo(HaveOccurred())
        })
    })
})
```

## See Also

- [[CI-CD|CI/CD]] — tests run as a required check on every PR
