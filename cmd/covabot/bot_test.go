package main

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCovaBot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CovaBot Suite")
}
