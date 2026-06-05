package engagement_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEngagement(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Engagement Suite")
}
