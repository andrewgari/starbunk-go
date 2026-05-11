package replybot_test

import (
	"path/filepath"
	"runtime"

	"github.com/andrewgari/starbunk-go/internal/replybot"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// testdataPath returns the absolute path to a file under testdata/.
func testdataPath(rel string) string {
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "testdata", rel)
}

var _ = Describe("Loader", func() {
	noMessaging := (*mockMessagingService)(nil)

	Describe("LoadFile", func() {
		It("loads all bots from a valid YAML file", func() {
			bots, err := replybot.LoadFile(testdataPath("single.yaml"), noMessaging)
			Expect(err).NotTo(HaveOccurred())
			Expect(bots).To(HaveLen(1))
			Expect(bots[0].Name).To(Equal("ping-bot"))
			Expect(bots[0].Triggers).To(HaveLen(1))
		})

		It("returns an error for a nonexistent file", func() {
			_, err := replybot.LoadFile(testdataPath("does-not-exist.yaml"), noMessaging)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error for invalid YAML", func() {
			_, err := replybot.LoadFile(testdataPath("invalid.yaml"), noMessaging)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error for a bot with an empty name", func() {
			_, err := replybot.LoadFile(testdataPath("bad_empty_name.yaml"), noMessaging)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("name must not be empty"))
		})

		It("returns an error for a bot with an invalid regex pattern", func() {
			_, err := replybot.LoadFile(testdataPath("bad_regex.yaml"), noMessaging)
			Expect(err).To(HaveOccurred())
		})

		It("returns an error for an unknown identity type", func() {
			_, err := replybot.LoadFile(testdataPath("bad_identity.yaml"), noMessaging)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unknown identity type"))
		})
	})

	Describe("LoadDir", func() {
		It("loads bots from all YAML files in a directory", func() {
			bots, err := replybot.LoadDir(testdataPath("multi"), noMessaging)
			Expect(err).NotTo(HaveOccurred())
			Expect(bots).To(HaveLen(2))
			// Alphabetical order: a.yaml then b.yaml
			Expect(bots[0].Name).To(Equal("hello-bot"))
			Expect(bots[1].Name).To(Equal("bye-bot"))
		})

		It("returns an empty slice for an empty directory", func() {
			bots, err := replybot.LoadDir(testdataPath("empty"), noMessaging)
			Expect(err).NotTo(HaveOccurred())
			Expect(bots).To(BeEmpty())
		})

		It("returns an error for a nonexistent directory", func() {
			_, err := replybot.LoadDir(testdataPath("no-such-dir"), noMessaging)
			Expect(err).To(HaveOccurred())
		})
	})
})
