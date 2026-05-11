package replybot_test

import (
	"github.com/andrewgari/starbunk-go/internal/replybot"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BuildIdentityResolver", func() {
	Describe("static identity", func() {
		It("returns a resolver that always yields the configured name and avatar", func() {
			resolver, err := replybot.BuildIdentityResolver(replybot.IdentityConfig{
				Type:      "static",
				BotName:   "Coolbot",
				AvatarURL: "https://example.com/avatar.png",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resolver).NotTo(BeNil())

			id, err := resolver(nil, buildMsg())
			Expect(err).NotTo(HaveOccurred())
			Expect(id.Username).To(Equal("Coolbot"))
			Expect(id.AvatarURL).To(Equal("https://example.com/avatar.png"))
		})

		It("errors when botName is empty", func() {
			_, err := replybot.BuildIdentityResolver(replybot.IdentityConfig{
				Type:      "static",
				AvatarURL: "https://example.com/avatar.png",
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("botName"))
		})

		It("errors when avatarUrl is empty", func() {
			_, err := replybot.BuildIdentityResolver(replybot.IdentityConfig{
				Type:    "static",
				BotName: "Coolbot",
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("avatarUrl"))
		})
	})

	Describe("mimic identity", func() {
		It("builds successfully when as_member is set", func() {
			resolver, err := replybot.BuildIdentityResolver(replybot.IdentityConfig{
				Type:     "mimic",
				AsMember: "123456789",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resolver).NotTo(BeNil())
		})

		It("errors when as_member is empty", func() {
			_, err := replybot.BuildIdentityResolver(replybot.IdentityConfig{
				Type: "mimic",
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("as_member"))
		})
	})

	Describe("random identity", func() {
		It("builds successfully without additional config", func() {
			resolver, err := replybot.BuildIdentityResolver(replybot.IdentityConfig{
				Type: "random",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resolver).NotTo(BeNil())
		})
	})

	Describe("unknown identity type", func() {
		It("returns an error", func() {
			_, err := replybot.BuildIdentityResolver(replybot.IdentityConfig{
				Type: "magic",
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unknown identity type"))
		})
	})
})
