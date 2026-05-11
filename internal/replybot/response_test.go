package replybot_test

import (
	"github.com/andrewgari/starbunk-go/internal/replybot"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// seededPool creates a ResponsePool with a deterministic rng (always picks index 0).
func seededPool(responses []string) *replybot.ResponsePool {
	p, err := replybot.NewResponsePoolWithRng(responses, func(_ int) int { return 0 })
	if err != nil {
		panic(err)
	}
	return p
}

var _ = Describe("ResponsePool", func() {
	Describe("NewResponsePool", func() {
		It("returns an error for an empty slice", func() {
			_, err := replybot.NewResponsePool([]string{})
			Expect(err).To(HaveOccurred())
		})
		It("succeeds with one or more responses", func() {
			p, err := replybot.NewResponsePool([]string{"hello"})
			Expect(err).NotTo(HaveOccurred())
			Expect(p).NotTo(BeNil())
		})
	})

	Describe("Pick", func() {
		It("returns the only item from a single-item pool", func() {
			p := seededPool([]string{"only response"})
			Expect(p.Pick("anything")).To(Equal("only response"))
		})

		It("returns the first item when rng always returns 0", func() {
			p := seededPool([]string{"first", "second", "third"})
			Expect(p.Pick("msg")).To(Equal("first"))
		})

		Describe("{start} template", func() {
			It("replaces {start} with the first three words", func() {
				p := seededPool([]string{"You said: {start}"})
				Expect(p.Pick("hello world foo bar baz")).To(Equal("You said: hello world foo"))
			})
			It("uses all words when fewer than three", func() {
				p := seededPool([]string{"{start}"})
				Expect(p.Pick("hi")).To(Equal("hi"))
				Expect(p.Pick("hi there")).To(Equal("hi there"))
			})
			It("returns empty string for empty message", func() {
				p := seededPool([]string{"{start}"})
				Expect(p.Pick("")).To(Equal(""))
			})
		})

		Describe("{random:min-max:char} template", func() {
			It("repeats the char exactly min times when rng returns 0", func() {
				p := seededPool([]string{"sh{random:2-5:e}sh"})
				// rng returns 0, so count = 2 + 0 = 2
				Expect(p.Pick("")).To(Equal("sheesh"))
			})
			It("handles min == max", func() {
				p := seededPool([]string{"{random:3-3:*}"})
				Expect(p.Pick("")).To(Equal("***"))
			})
			It("leaves malformed placeholder unchanged", func() {
				p := seededPool([]string{"{random:bad-x:e}"})
				Expect(p.Pick("")).To(Equal("{random:bad-x:e}"))
			})
		})

		Describe("{swap_message:word1:word2} template", func() {
			It("swaps word1 for word2 in the original message", func() {
				p := seededPool([]string{"{swap_message:check:czech}"})
				Expect(p.Pick("check this out")).To(Equal("czech this out"))
			})
			It("is case-insensitive on the target word", func() {
				p := seededPool([]string{"{swap_message:foo:bar}"})
				Expect(p.Pick("FOO and foo")).To(Equal("bar and bar"))
			})
			It("returns the original message when the word is not found", func() {
				p := seededPool([]string{"{swap_message:xyz:abc}"})
				Expect(p.Pick("nothing to swap here")).To(Equal("nothing to swap here"))
			})
		})

		Describe("multiple templates in one response", func() {
			It("expands all templates independently", func() {
				p := seededPool([]string{"swap: {swap_message:a:b} | start: {start}"})
				Expect(p.Pick("a b c d")).To(Equal("swap: b b c d | start: a b c"))
			})
		})
	})
})
