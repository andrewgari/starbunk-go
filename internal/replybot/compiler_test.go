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
func withChannelID(id string) msgOpt {
	return func(m *discordgo.MessageCreate) { m.ChannelID = id }
}

// intPtr is a helper to take the address of an int literal.
func intPtr(n int) *int { return &n }

// — compiler tests —

var _ = Describe("Compile", func() {
	s := testSession()

	Describe("leaf conditions — content", func() {
		Context("contains_word", func() {
			It("matches a whole word case-insensitively", func() {
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

		Context("always", func() {
			It("passes unconditionally", func() {
				a, err := replybot.Compile(replybot.ConditionNode{Always: true})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg())).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("anything")))).To(BeTrue())
			})
		})
	})

	Describe("leaf conditions — author", func() {
		Context("from_user", func() {
			It("passes for the matching user ID only", func() {
				a, err := replybot.Compile(replybot.ConditionNode{FromUser: "u123"})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withAuthorID("u123")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withAuthorID("u456")))).To(BeFalse())
			})
		})

		Context("not_from_user", func() {
			It("passes for any user except the excluded one", func() {
				a, err := replybot.Compile(replybot.ConditionNode{NotFromUser: "u123"})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withAuthorID("u456")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withAuthorID("u123")))).To(BeFalse())
			})
		})

		Context("author_is_bot", func() {
			It("passes for bot authors, fails for humans", func() {
				a, err := replybot.Compile(replybot.ConditionNode{AuthorIsBot: true})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withBot()))).To(BeTrue())
				Expect(a.Audit(s, buildMsg())).To(BeFalse())
			})
		})

		Context("author_not_bot", func() {
			It("passes for humans, fails for bots", func() {
				a, err := replybot.Compile(replybot.ConditionNode{AuthorNotBot: true})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg())).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withBot()))).To(BeFalse())
			})
		})
	})

	Describe("leaf conditions — context", func() {
		Context("in_channel", func() {
			It("passes only for the specified channel", func() {
				a, err := replybot.Compile(replybot.ConditionNode{InChannel: "chan-1"})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withChannelID("chan-1")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withChannelID("chan-2")))).To(BeFalse())
			})
		})
	})

	Describe("leaf conditions — probabilistic", func() {
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
	})

	Describe("error cases", func() {
		It("returns an error for an empty node", func() {
			_, err := replybot.Compile(replybot.ConditionNode{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("empty condition node"))
		})
		It("returns an error for empty combinator slices", func() {
			_, err := replybot.Compile(replybot.ConditionNode{AllOf: []replybot.ConditionNode{}})
			Expect(err).To(HaveOccurred())
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
				Expect(a.Audit(s, buildMsg(withContent("only foo here")))).To(BeFalse())
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
				Expect(a.Audit(s, buildMsg(withContent("this is spam")))).To(BeFalse())
			})
		})
	})

	Describe("combination scenarios", func() {
		Context("scenario 1 — bot detection (author_is_bot alone)", func() {
			It("fires only on bot messages", func() {
				a, err := replybot.Compile(replybot.ConditionNode{AuthorIsBot: true})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withBot()))).To(BeTrue())
				Expect(a.Audit(s, buildMsg())).To(BeFalse())
			})
		})

		Context("scenario 1b — bot AND specific phrase", func() {
			It("fires only when a bot says a specific word", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					AllOf: []replybot.ConditionNode{
						{AuthorIsBot: true},
						{ContainsPhrase: "hello"},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withBot(), withContent("hello world")))).To(BeTrue())
				// human saying hello — no match
				Expect(a.Audit(s, buildMsg(withContent("hello world")))).To(BeFalse())
				// bot saying something else — no match
				Expect(a.Audit(s, buildMsg(withBot(), withContent("goodbye")))).To(BeFalse())
			})
		})

		Context("scenario 2 — specific user + specific phrase", func() {
			It("fires only when user u42 says the trigger phrase", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					AllOf: []replybot.ConditionNode{
						{FromUser: "u42"},
						{ContainsPhrase: "open sesame"},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				// Correct user, correct phrase
				Expect(a.Audit(s, buildMsg(withAuthorID("u42"), withContent("open sesame please")))).To(BeTrue())
				// Wrong user
				Expect(a.Audit(s, buildMsg(withAuthorID("u99"), withContent("open sesame please")))).To(BeFalse())
				// Correct user, wrong phrase
				Expect(a.Audit(s, buildMsg(withAuthorID("u42"), withContent("close sesame")))).To(BeFalse())
			})
		})

		Context("scenario 4 — regex pattern matching", func() {
			It("fires on messages containing a phone-number-like pattern", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					MatchesPattern: `\d{3}[-.\s]\d{3}[-.\s]\d{4}`,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("call me at 555-867-5309")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("no number here")))).To(BeFalse())
			})
		})

		Context("regex + from_user (specific user using a pattern)", func() {
			It("fires only when that user sends a message matching the regex", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					AllOf: []replybot.ConditionNode{
						{FromUser: "u42"},
						{MatchesPattern: `(?i)^!roll\s+\d+d\d+`},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withAuthorID("u42"), withContent("!roll 2d6")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withAuthorID("u99"), withContent("!roll 2d6")))).To(BeFalse())
				Expect(a.Audit(s, buildMsg(withAuthorID("u42"), withContent("just talking")))).To(BeFalse())
			})
		})

		Context("any_of + multiple keywords (greeting detector)", func() {
			It("fires on any of several greeting words", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					AnyOf: []replybot.ConditionNode{
						{ContainsWord: "hello"},
						{ContainsWord: "hi"},
						{ContainsWord: "howdy"},
						{ContainsPhrase: "hey there"},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("hello everyone")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("hi there")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("howdy partner")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("hey there mate")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("greetings")))).To(BeFalse())
			})
		})

		Context("none_of to exclude bots and a specific user", func() {
			It("fires for everyone except bots and the excluded user", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					AllOf: []replybot.ConditionNode{
						{ContainsWord: "test"},
						{
							NoneOf: []replybot.ConditionNode{
								{AuthorIsBot: true},
								{FromUser: "excluded-u"},
							},
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				// Normal human saying test — passes
				Expect(a.Audit(s, buildMsg(withAuthorID("u1"), withContent("test")))).To(BeTrue())
				// Bot saying test — blocked by none_of
				Expect(a.Audit(s, buildMsg(withBot(), withContent("test")))).To(BeFalse())
				// Excluded user saying test — blocked by none_of
				Expect(a.Audit(s, buildMsg(withAuthorID("excluded-u"), withContent("test")))).To(BeFalse())
				// Normal human not saying test — blocked by first all_of
				Expect(a.Audit(s, buildMsg(withAuthorID("u1"), withContent("hello")))).To(BeFalse())
			})
		})

		Context("all_of + any_of nested (two groups must both match)", func() {
			It("fires when message has a fruit AND a drink", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					AllOf: []replybot.ConditionNode{
						{
							AnyOf: []replybot.ConditionNode{
								{ContainsWord: "banana"},
								{ContainsWord: "apple"},
							},
						},
						{
							AnyOf: []replybot.ConditionNode{
								{ContainsWord: "juice"},
								{ContainsWord: "smoothie"},
							},
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("banana juice")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("apple smoothie")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("banana")))).To(BeFalse())
				Expect(a.Audit(s, buildMsg(withContent("juice only")))).To(BeFalse())
			})
		})

		Context("word + chance (50% random response to a keyword)", func() {
			It("fires on keyword at 100% chance", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					AllOf: []replybot.ConditionNode{
						{ContainsWord: "fish"},
						{WithChance: intPtr(100)},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withContent("I like fish")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withContent("no keyword")))).To(BeFalse())
			})
			It("never fires on keyword at 0% chance", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					AllOf: []replybot.ConditionNode{
						{ContainsWord: "fish"},
						{WithChance: intPtr(0)},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				for range 5 {
					Expect(a.Audit(s, buildMsg(withContent("I like fish")))).To(BeFalse())
				}
			})
		})

		Context("in_channel + regex (channel-scoped pattern match)", func() {
			It("fires only in the target channel with a matching pattern", func() {
				a, err := replybot.Compile(replybot.ConditionNode{
					AllOf: []replybot.ConditionNode{
						{InChannel: "announcements"},
						{MatchesPattern: `(?i)@everyone`},
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(a.Audit(s, buildMsg(withChannelID("announcements"), withContent("@everyone listen up!")))).To(BeTrue())
				Expect(a.Audit(s, buildMsg(withChannelID("general"), withContent("@everyone listen up!")))).To(BeFalse())
				Expect(a.Audit(s, buildMsg(withChannelID("announcements"), withContent("just talking")))).To(BeFalse())
			})
		})
	})
})
