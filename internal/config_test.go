package pkg_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestInternal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Internal Suite")
}

var _ = Describe("Config", func() {
	It("should compile the test correctly", func() {
		// Just a basic test to verify Ginkgo works
		Expect(true).To(BeTrue())
	})
})
