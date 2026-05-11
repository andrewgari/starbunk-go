package replybot_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestReplyBot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ReplyBot Suite")
}
