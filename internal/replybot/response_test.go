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

// msgData is a convenience constructor for MessageData in tests.
func msgData(content, authorUsername, authorID string) replybot.MessageData {
	return replybot.MessageData{
		Content:        content,
		AuthorUsername: authorUsername,
		AuthorID:       authorID,
	}
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
			Expect(p.Pick(msgData("anything", "user", "u1"))).To(Equal("only response"))
		})

		It("returns the first item when rng always returns 0", func() {
			p := seededPool([]string{"first", "second", "third"})
			Expect(p.Pick(msgData("msg", "user", "u1"))).To(Equal("first"))
		})

		Describe("{start} template", func() {
			It("replaces {start} with the first three words", func() {
				p := seededPool([]string{"You said: {start}"})
				Expect(p.Pick(msgData("hello world foo bar baz", "user", "u1"))).To(Equal("You said: hello world foo"))
			})
			It("uses all words when fewer than three", func() {
				p := seededPool([]string{"{start}"})
				Expect(p.Pick(msgData("hi", "user", "u1"))).To(Equal("hi"))
				Expect(p.Pick(msgData("hi there", "user", "u1"))).To(Equal("hi there"))
			})
			It("returns empty string for empty message", func() {
				p := seededPool([]string{"{start}"})
				Expect(p.Pick(msgData("", "user", "u1"))).To(Equal(""))
			})
		})

		Describe("{message} template", func() {
			It("expands to the full message content", func() {
				p := seededPool([]string{"Full message: {message}"})
				Expect(p.Pick(msgData("hello world foo bar baz", "user", "u1"))).To(Equal("Full message: hello world foo bar baz"))
			})
			It("returns empty string for empty content", func() {
				p := seededPool([]string{"{message}"})
				Expect(p.Pick(msgData("", "user", "u1"))).To(Equal(""))
			})
		})

		Describe("{author} template", func() {
			It("expands to the author's username", func() {
				p := seededPool([]string{"Hello, {author}!"})
				Expect(p.Pick(msgData("hi", "andrewgari", "u1"))).To(Equal("Hello, andrewgari!"))
			})
		})

		Describe("{author_id} template", func() {
			It("expands to the author's Discord ID", func() {
				p := seededPool([]string{"ID: {author_id}"})
				Expect(p.Pick(msgData("hi", "andrewgari", "123456789"))).To(Equal("ID: 123456789"))
			})
		})

		Describe("{author_mention} template", func() {
			It("expands to a Discord @mention string", func() {
				p := seededPool([]string{"Hey {author_mention}!"})
				Expect(p.Pick(msgData("hi", "andrewgari", "123456789"))).To(Equal("Hey <@123456789>!"))
			})
		})

		Describe("{random:min-max:char} template", func() {
			It("repeats the char exactly min times when rng returns 0", func() {
				p := seededPool([]string{"sh{random:2-5:e}sh"})
				// rng returns 0, count = 2 + 0 = 2
				Expect(p.Pick(msgData("", "user", "u1"))).To(Equal("sheesh"))
			})
			It("handles min == max", func() {
				p := seededPool([]string{"{random:3-3:*}"})
				Expect(p.Pick(msgData("", "user", "u1"))).To(Equal("***"))
			})
			It("leaves malformed placeholder unchanged", func() {
				p := seededPool([]string{"{random:bad-x:e}"})
				Expect(p.Pick(msgData("", "user", "u1"))).To(Equal("{random:bad-x:e}"))
			})
		})

		Describe("{swap_message:word1:word2} template", func() {
			It("swaps word1 for word2 in the original message", func() {
				p := seededPool([]string{"{swap_message:check:czech}"})
				Expect(p.Pick(msgData("check this out", "user", "u1"))).To(Equal("czech this out"))
			})
			It("is case-insensitive on the target word", func() {
				p := seededPool([]string{"{swap_message:foo:bar}"})
				Expect(p.Pick(msgData("FOO and foo", "user", "u1"))).To(Equal("bar and bar"))
			})
			It("returns the original message when the word is not found", func() {
				p := seededPool([]string{"{swap_message:xyz:abc}"})
				Expect(p.Pick(msgData("nothing to swap here", "user", "u1"))).To(Equal("nothing to swap here"))
			})
		})

		Describe("multiple templates in one response", func() {
			It("expands all templates independently", func() {
				p := seededPool([]string{"{author} said: {message} | start: {start}"})
				Expect(p.Pick(msgData("a b c d", "andrewgari", "u1"))).To(Equal("andrewgari said: a b c d | start: a b c"))
			})
			It("combines author_mention with start", func() {
				p := seededPool([]string{"Hey {author_mention}, you said {start}!"})
				Expect(p.Pick(msgData("hello world foo", "andrewgari", "999"))).To(Equal("Hey <@999>, you said hello world foo!"))
			})
		})
	})
})
