package replybot_test

import (
	"time"

	"github.com/andrewgari/starbunk-go/internal/replybot"
	"github.com/bwmarrin/discordgo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// — test helpers (shared across *_test.go files in this package) —

func testSession() *discordgo.Session {
	s := &discordgo.Session{State: discordgo.NewState()}
	s.State.User = &discordgo.User{ID: "bot-id"}
	return s
}

type msgOpt func(*discordgo.MessageCreate)

func buildMsg(opts ...msgOpt) *discordgo.MessageCreate {
	m := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author:    &discordgo.User{},
			Timestamp: time.Now(),
		},
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

func withContent(c string) msgOpt {
	return func(m *discordgo.MessageCreate) { m.Content = c }
}
func withAuthorID(id string) msgOpt {
	return func(m *discordgo.MessageCreate) { m.Author.ID = id }
}
func withBot() msgOpt {
	return func(m *discordgo.MessageCreate) { m.Author.Bot = true }
}
func withGuildID(id string) msgOpt {
	return func(m *discordgo.MessageCreate) { m.GuildID = id }
}

// intPtr is a helper to take the address of an int literal.
func intPtr(n int) *int { return &n }

// — compiler tests —

var _ = Describe("Compile", func() {
	s := testSession()

	Describe("leaf conditions", func() {
		Context("contains_word", func() {
			It("matches a whole word (case-insensitive)", func() {
				a, err := replybot.Compile(replybot.ConditionNode{ContainsWord: "hello"})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("Hello world")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("say hello to me")))).To(BeTrue())
			})
			It("does not match partial words", func() {
				a, err := replybot.Compile(replybot.ConditionNode{ContainsWord: "hello"})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("helloWorld")))).To(BeFalse())
				Expect(a.Audit(s, buildMsg(withContent("unhello")))).To(BeFalse())
			})
		})

		Context("contains_phrase", func() {
			It("matches a case-insensitive substring", func() {
				a, err := replybot.Compile(replybot.ConditionNode{ContainsPhrase: "blue mage"})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("I am a Blue Mage")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("blue mage is cool")))).To(BeTrue())
			})
			It("does not match when phrase is absent", func() {
				a, err := replybot.Compile(replybot.ConditionNode{ContainsPhrase: "blue mage"})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("red wizard")))).To(BeFalse())
			})
		})

		Context("matches_pattern", func() {
			It("matches a valid regex", func() {
				a, err := replybot.Compile(replybot.ConditionNode{MatchesPattern: `\d+`})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("call 911")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("no numbers here")))).To(BeFalse())
			})
			It("returns an error for invalid regex", func() {
				_, err := replybot.Compile(replybot.ConditionNode{MatchesPattern: "["})
				Expect(err).To(HaveOccurred())
			})
		})

		Context("from_user", func() {
			It("passes for the matching user ID", func() {
				a, err := replybot.Compile(replybot.ConditionNode{FromUser: "u123"})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withAuthorID("u123")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withAuthorID("u456")))).To(BeFalse())
			})
		})

		Context("with_chance", func() {
			It("always passes at 100", func() {
				a, err := replybot.Compile(replybot.ConditionNode{WithChance: intPtr(100)})
				Expect(err).NotTo(HaveOccurred())
				for range 10 {
					Expect(a.Audit(s, buildMsg())).To(BeTrue())
				}
			})
			It("never passes at 0", func() {
				a, err := replybot.Compile(replybot.ConditionNode{WithChance: intPtr(0)})
				Expect(err).NotTo(HaveOccurred())
				for range 10 {
					Expect(a.Audit(s, buildMsg())).To(BeFalse())
				}
			})
		})

		Context("always", func() {
			It("passes unconditionally", func() {
				a, err := replybot.Compile(replybot.ConditionNode{Always: true})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg())).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("anything")))).To(BeTrue())
			})
		})

		It("returns an error for an empty node", func() {
			_, err := replybot.Compile(replybot.ConditionNode{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("empty condition node"))
		})
	})

	Describe("logical combinators", func() {
		Context("all_of", func() {
			It("passes when all children pass", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					AllOf: []replybot.ConditionNode{
						{ContainsWord: "foo"},
						{ContainsWord: "bar"},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("foo and bar")))).To(BeTrue())
			})
			It("fails when any child fails", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					AllOf: []replybot.ConditionNode{
						{ContainsWord: "foo"},
						{ContainsWord: "bar"},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("only foo here")))).To(BeFalse())
			})
			It("returns an error for empty all_of", func() {
				_, err := replybot.Compile(replybot.ConditionNode{AllOf: []replybot.ConditionNode{}})
				Expect(err).To(HaveOccurred())
			})
		})

		Context("any_of", func() {
			It("passes when at least one child passes", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					AnyOf: []replybot.ConditionNode{
						{ContainsWord: "hello"},
						{ContainsWord: "hi"},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("hi there")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("hello world")))).To(BeTrue())
			})
			It("fails when no child passes", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					AnyOf: []replybot.ConditionNode{
						{ContainsWord: "hello"},
						{ContainsWord: "hi"},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("greetings")))).To(BeFalse())
			})
		})

		Context("none_of", func() {
			It("passes when no child passes", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					NoneOf: []replybot.ConditionNode{
						{ContainsWord: "spam"},
						{ContainsWord: "junk"},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("clean message")))).To(BeTrue())
			})
			It("fails when any child passes", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					NoneOf: []replybot.ConditionNode{
						{ContainsWord: "spam"},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("this is spam")))).To(BeFalse())
			})
		})

		Context("nested logic", func() {
			It("handles all_of containing any_of", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					AllOf: []replybot.ConditionNode{
						{
							AnyOf: []replybot.ConditionNode{
								{ContainsWord: "banana"},
								{ContainsWord: "apple"},
							},
						},
						{ContainsWord: "juice"},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("banana juice")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("apple juice")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("orange juice")))).To(BeFalse())
				Expect(a.Audit(s, buildMsg(withContent("banana smoothie")))).To(BeFalse())
			})
		})
	})
})
