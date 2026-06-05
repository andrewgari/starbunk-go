package tagger_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTagger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tagger Suite")
}
