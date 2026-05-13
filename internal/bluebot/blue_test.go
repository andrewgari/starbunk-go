package bluebot_test

import (
	"context"
	"testing"
	"time"

	"github.com/andrewgari/starbunk-go/internal/bluebot"
	"github.com/bwmarrin/discordgo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBluebot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bluebot Suite")
}

func msg(content string) *discordgo.MessageCreate {
	return &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Content:   content,
			Author:    &discordgo.User{},
			Timestamp: time.Now(),
		},
	}
}

var _ = Describe("BlueStrategy", func() {
	var (
		s   bluebot.BlueStrategy
		ctx = context.Background()
	)

	Describe("Name", func() {
		It("identifies itself", func() {
			Expect(s.Name()).To(Equal("BlueStrategy"))
		})
	})

	Describe("ShouldTrigger", func() {
		DescribeTable("matches references to blue",
			func(content string) {
				Expect(s.ShouldTrigger(ctx, msg(content))).To(BeTrue())
			},
			Entry("plain lowercase", "i like blue"),
			Entry("title case", "Blue is my favourite colour"),
			Entry("all caps", "BLUE"),
			Entry("Blue Mage job reference", "I play Blue Mage"),
			Entry("mid-sentence", "the sky is blue today"),
			Entry("blu short form", "blu"),
			Entry("bluu elongated", "bluu"),
			Entry("bloo elongated o", "bloo"),
			Entry("blooo very elongated", "blooo"),
			Entry("blew homophone", "he blew the whistle"),
			Entry("bleu French spelling", "cordon bleu"),
			Entry("azul Spanish/Portuguese", "azul"),
			Entry("blau German", "das ist blau"),
			Entry("bluebot self-reference", "hey bluebot"),
			Entry("punctuation after", "blue!"),
			Entry("punctuation before", "say: blue"),
		)

		DescribeTable("does not match unrelated content",
			func(content string) {
				Expect(s.ShouldTrigger(ctx, msg(content))).To(BeFalse())
			},
			Entry("unrelated colour", "I like red"),
			Entry("bluetooth — compound word", "connect via bluetooth"),
			Entry("blueprint — compound word", "read the blueprint"),
			Entry("blueberry — compound word", "eat a blueberry"),
			Entry("empty message", ""),
			Entry("numbers only", "12345"),
			Entry("unrelated word", "hello world"),
		)
	})

	Describe("Response", func() {
		It("returns the classic catchphrase", func() {
			Expect(s.Response(ctx, msg("blue"))).To(Equal("Did somebody say Blu?"))
		})
	})
})
