package main

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBunkBot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BunkBot Suite")
}
