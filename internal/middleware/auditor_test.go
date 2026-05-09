package middleware_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/andrewgari/starbunk-go/internal/middleware"
	"github.com/bwmarrin/discordgo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Middleware Suite")
}

// — test helpers —

func session(botID string) *discordgo.Session {
	s := &discordgo.Session{State: discordgo.NewState()}
	s.State.User = &discordgo.User{ID: botID}
	return s
}

type msgOpt func(*discordgo.MessageCreate)

func build(opts ...msgOpt) *discordgo.MessageCreate {
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

func authorID(id string) msgOpt {
	return func(m *discordgo.MessageCreate) { m.Author.ID = id }
}
func authorName(name string) msgOpt {
	return func(m *discordgo.MessageCreate) { m.Author.Username = name }
}
func isBot() msgOpt {
	return func(m *discordgo.MessageCreate) { m.Author.Bot = true }
}
func content(c string) msgOpt {
	return func(m *discordgo.MessageCreate) { m.Content = c }
}
func guildID(id string) msgOpt {
	return func(m *discordgo.MessageCreate) { m.GuildID = id }
}
func channelID(id string) msgOpt {
	return func(m *discordgo.MessageCreate) { m.ChannelID = id }
}
func attachment() msgOpt {
	return func(m *discordgo.MessageCreate) {
		m.Attachments = []*discordgo.MessageAttachment{{ID: "att-1"}}
	}
}
func sentOn(day time.Weekday) msgOpt {
	// Advance from a known Monday until we reach the target weekday.
	monday := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC) // 2024-01-01 is a Monday
	offset := (int(day) - int(time.Monday) + 7) % 7
	ts := monday.AddDate(0, 0, offset)
	return func(m *discordgo.MessageCreate) { m.Timestamp = ts }
}

// — author gates —

var _ = Describe("NotSelf", func() {
	s := session("bot-id")

	It("passes messages from other users", func() {
		Expect(middleware.NotSelf.Audit(s, build(authorID("user-id")))).To(BeTrue())
	})
	It("drops messages from the bot itself", func() {
		Expect(middleware.NotSelf.Audit(s, build(authorID("bot-id")))).To(BeFalse())
	})
	It("passes when State.User is nil (session not ready)", func() {
		bare := &discordgo.Session{State: discordgo.NewState()}
		Expect(middleware.NotSelf.Audit(bare, build(authorID("anyone")))).To(BeTrue())
	})
})

var _ = Describe("NotBot / IsBot", func() {
	s := session("bot-id")

	It("NotBot passes human messages", func() {
		Expect(middleware.NotBot.Audit(s, build(authorID("user")))).To(BeTrue())
	})
	It("NotBot drops bot messages", func() {
		Expect(middleware.NotBot.Audit(s, build(authorID("other-bot"), isBot()))).To(BeFalse())
	})
	It("IsBot passes bot messages", func() {
		Expect(middleware.IsBot.Audit(s, build(authorID("other-bot"), isBot()))).To(BeTrue())
	})
	It("IsBot drops human messages", func() {
		Expect(middleware.IsBot.Audit(s, build(authorID("user")))).To(BeFalse())
	})
})

var _ = Describe("AuthorID", func() {
	s := session("bot-id")

	It("passes the matching author", func() {
		Expect(middleware.AuthorID("u1").Audit(s, build(authorID("u1")))).To(BeTrue())
	})
	It("drops a different author", func() {
		Expect(middleware.AuthorID("u1").Audit(s, build(authorID("u2")))).To(BeFalse())
	})
})

var _ = Describe("NotAuthorID", func() {
	s := session("bot-id")

	It("drops the specified author", func() {
		Expect(middleware.NotAuthorID("u1").Audit(s, build(authorID("u1")))).To(BeFalse())
	})
	It("passes any other author", func() {
		Expect(middleware.NotAuthorID("u1").Audit(s, build(authorID("u2")))).To(BeTrue())
	})
})

var _ = Describe("AuthorNamed", func() {
	s := session("bot-id")

	It("passes when username matches", func() {
		Expect(middleware.AuthorNamed("Jeff").Audit(s, build(authorName("Jeff")))).To(BeTrue())
	})
	It("drops when username does not match", func() {
		Expect(middleware.AuthorNamed("Jeff").Audit(s, build(authorName("Bob")))).To(BeFalse())
	})
	It("is case-sensitive", func() {
		Expect(middleware.AuthorNamed("Jeff").Audit(s, build(authorName("jeff")))).To(BeFalse())
	})
})

// — content gates —

var _ = Describe("HasContent", func() {
	s := session("bot-id")

	It("passes messages with text", func() {
		Expect(middleware.HasContent.Audit(s, build(content("hello")))).To(BeTrue())
	})
	It("drops empty content", func() {
		Expect(middleware.HasContent.Audit(s, build(content("")))).To(BeFalse())
	})
	It("drops whitespace-only content", func() {
		Expect(middleware.HasContent.Audit(s, build(content("  \t\n")))).To(BeFalse())
	})
})

var _ = Describe("ContentContains", func() {
	s := session("bot-id")

	It("passes when the substring is present", func() {
		Expect(middleware.ContentContains("blue").Audit(s, build(content("I love blue")))).To(BeTrue())
	})
	It("drops when the substring is absent", func() {
		Expect(middleware.ContentContains("blue").Audit(s, build(content("I love red")))).To(BeFalse())
	})
})

var _ = Describe("ContentMatches", func() {
	s := session("bot-id")
	re := regexp.MustCompile(`(?i)\bblue\b`)

	It("passes when the pattern matches", func() {
		Expect(middleware.ContentMatches(re).Audit(s, build(content("Blue Mage is cool")))).To(BeTrue())
	})
	It("drops when the pattern does not match", func() {
		Expect(middleware.ContentMatches(re).Audit(s, build(content("Red Mage is cool")))).To(BeFalse())
	})
})

var _ = Describe("HasAttachment", func() {
	s := session("bot-id")

	It("passes when an attachment is present", func() {
		Expect(middleware.HasAttachment.Audit(s, build(attachment()))).To(BeTrue())
	})
	It("drops when no attachments", func() {
		Expect(middleware.HasAttachment.Audit(s, build())).To(BeFalse())
	})
})

// — context gates —

var _ = Describe("GuildOnly", func() {
	s := session("bot-id")

	It("passes guild messages", func() {
		Expect(middleware.GuildOnly.Audit(s, build(guildID("g")))).To(BeTrue())
	})
	It("drops DMs", func() {
		Expect(middleware.GuildOnly.Audit(s, build())).To(BeFalse())
	})
})

var _ = Describe("DMOnly", func() {
	s := session("bot-id")

	It("passes DMs", func() {
		Expect(middleware.DMOnly.Audit(s, build())).To(BeTrue())
	})
	It("drops guild messages", func() {
		Expect(middleware.DMOnly.Audit(s, build(guildID("g")))).To(BeFalse())
	})
})

var _ = Describe("InChannel", func() {
	s := session("bot-id")

	It("passes messages in the specified channel", func() {
		Expect(middleware.InChannel("chan-1").Audit(s, build(channelID("chan-1")))).To(BeTrue())
	})
	It("drops messages in a different channel", func() {
		Expect(middleware.InChannel("chan-1").Audit(s, build(channelID("chan-2")))).To(BeFalse())
	})
})

var _ = Describe("OnWeekdays", func() {
	s := session("bot-id")

	It("passes on a listed weekday", func() {
		Expect(middleware.OnWeekdays(time.Monday, time.Wednesday).Audit(s, build(sentOn(time.Monday)))).To(BeTrue())
	})
	It("drops on an unlisted weekday", func() {
		Expect(middleware.OnWeekdays(time.Monday, time.Wednesday).Audit(s, build(sentOn(time.Friday)))).To(BeFalse())
	})
})

// — combinators —

var _ = Describe("AllOf", func() {
	s := session("bot-id")

	It("passes when all children pass", func() {
		a := middleware.AllOf(middleware.NotSelf, middleware.HasContent)
		Expect(a.Audit(s, build(authorID("u"), content("hi")))).To(BeTrue())
	})
	It("drops when the first child fails", func() {
		a := middleware.AllOf(middleware.NotSelf, middleware.HasContent)
		Expect(a.Audit(s, build(authorID("bot-id"), content("hi")))).To(BeFalse())
	})
	It("drops when a later child fails", func() {
		a := middleware.AllOf(middleware.NotSelf, middleware.HasContent)
		Expect(a.Audit(s, build(authorID("u"), content("")))).To(BeFalse())
	})
	It("passes vacuously with no children", func() {
		Expect(middleware.AllOf().Audit(s, build())).To(BeTrue())
	})
})

var _ = Describe("AnyOf", func() {
	s := session("bot-id")

	It("passes when the first child passes", func() {
		Expect(middleware.AnyOf(middleware.GuildOnly, middleware.DMOnly).Audit(s, build(guildID("g")))).To(BeTrue())
	})
	It("passes when only a later child passes", func() {
		Expect(middleware.AnyOf(middleware.GuildOnly, middleware.DMOnly).Audit(s, build())).To(BeTrue())
	})
	It("drops when no child passes", func() {
		neverPass := middleware.AllOf(middleware.GuildOnly, middleware.DMOnly)
		Expect(middleware.AnyOf(neverPass, neverPass).Audit(s, build(guildID("g")))).To(BeFalse())
	})
	It("fails vacuously with no children", func() {
		Expect(middleware.AnyOf().Audit(s, build())).To(BeFalse())
	})
})

var _ = Describe("Not", func() {
	s := session("bot-id")

	It("inverts a passing auditor", func() {
		Expect(middleware.Not(middleware.GuildOnly).Audit(s, build(guildID("g")))).To(BeFalse())
	})
	It("inverts a failing auditor", func() {
		Expect(middleware.Not(middleware.GuildOnly).Audit(s, build())).To(BeTrue())
	})
})

// — Chance —

var _ = Describe("Chance", func() {
	s := session("bot-id")

	It("always passes with probability 1.0", func() {
		// rand.Float64 returns [0.0, 1.0), so it is always < 1.0
		Expect(middleware.Chance(1.0).Audit(s, build())).To(BeTrue())
	})
	It("never passes with probability 0.0", func() {
		// rand.Float64 returns [0.0, 1.0), so it is never < 0.0
		Expect(middleware.Chance(0.0).Audit(s, build())).To(BeFalse())
	})
})

// — complex compositions —

var _ = Describe("complex compositions", func() {
	s := session("bot-id")

	It("drops a bot named Jeff messaging on Friday or Tuesday", func() {
		// Represents an auditor that REJECTS that specific scenario.
		// Read: reject if (IsBot AND AuthorNamed("Jeff") AND OnWeekdays(Friday, Tuesday))
		rejectJeffOnWeekends := middleware.Not(middleware.AllOf(
			middleware.IsBot,
			middleware.AuthorNamed("Jeff"),
			middleware.OnWeekdays(time.Friday, time.Tuesday),
		))

		jeffBotOnFriday := build(isBot(), authorName("Jeff"), sentOn(time.Friday))
		jeffBotOnMonday := build(isBot(), authorName("Jeff"), sentOn(time.Monday))
		humanJeffOnFriday := build(authorName("Jeff"), sentOn(time.Friday))

		Expect(rejectJeffOnWeekends.Audit(s, jeffBotOnFriday)).To(BeFalse()) // blocked
		Expect(rejectJeffOnWeekends.Audit(s, jeffBotOnMonday)).To(BeTrue())  // Jeff, but safe day
		Expect(rejectJeffOnWeekends.Audit(s, humanJeffOnFriday)).To(BeTrue()) // human Jeff, passes
	})

	It("BlueBot policy: passes human guild messages with content", func() {
		bluebotAuditor := middleware.AllOf(
			middleware.NotSelf,
			middleware.NotBot,
			middleware.GuildOnly,
			middleware.HasContent,
		)

		Expect(bluebotAuditor.Audit(s, build(authorID("u"), guildID("g"), content("hello")))).To(BeTrue())
		Expect(bluebotAuditor.Audit(s, build(authorID("bot-id"), guildID("g"), content("hello")))).To(BeFalse()) // self
		Expect(bluebotAuditor.Audit(s, build(authorID("u"), isBot(), guildID("g"), content("hi")))).To(BeFalse()) // bot author
		Expect(bluebotAuditor.Audit(s, build(authorID("u"), content("hello")))).To(BeFalse())                    // DM
		Expect(bluebotAuditor.Audit(s, build(authorID("u"), guildID("g"), content("")))).To(BeFalse())           // empty
	})

	It("BunkBot policy: passes non-self guild messages with content, including from bots", func() {
		bunkbotAuditor := middleware.AllOf(
			middleware.NotSelf,
			middleware.HasContent,
		)

		Expect(bunkbotAuditor.Audit(s, build(authorID("other-bot"), isBot(), guildID("g"), content("hi")))).To(BeTrue())
		Expect(bunkbotAuditor.Audit(s, build(authorID("bot-id"), guildID("g"), content("hi")))).To(BeFalse()) // self
		Expect(bunkbotAuditor.Audit(s, build(authorID("u"), guildID("g"), content("")))).To(BeFalse())        // empty
	})

	// Scenario 1:
	// Bots always fail (hard gate).
	// Among non-bots: userid 111111 only triggers when content contains "bingo".
	// All other non-bot users trigger freely.
	//
	//   AllOf(
	//     NotBot,
	//     AnyOf(Not(AuthorID("111111")), ContentContains("bingo")),
	//   )
	It("scenario 1: bots always fail; non-bots pass freely except userid 111111 who needs 'bingo'", func() {
		auditor := middleware.AllOf(
			middleware.NotBot,
			middleware.AnyOf(
				middleware.Not(middleware.AuthorID("111111")),
				middleware.ContentContains("bingo"),
			),
		)

		// Any bot → fails, regardless of user ID or content
		Expect(auditor.Audit(s, build(isBot(), authorID("111111"), content("bingo")))).To(BeFalse())
		Expect(auditor.Audit(s, build(isBot(), authorID("999999"), content("anything")))).To(BeFalse())
		// Non-bot, not userid 111111 → passes freely
		Expect(auditor.Audit(s, build(authorID("999999"), content("anything")))).To(BeTrue())
		// Non-bot, userid 111111, content contains "bingo" → passes
		Expect(auditor.Audit(s, build(authorID("111111"), content("bingo")))).To(BeTrue())
		// Non-bot, userid 111111, no "bingo" → fails
		Expect(auditor.Audit(s, build(authorID("111111"), content("hello")))).To(BeFalse())
	})

	// Scenario 2:
	// Accept a message if: NOT a bot
	//                      OR  (author is 22222 AND 1% chance)
	//                      OR  content contains "bingo"
	It("scenario 2: human passes; 'bingo' always passes; bot 22222 passes on a lucky roll", func() {
		// Use Chance(1.0) to represent a winning roll and Chance(0.0) for a losing roll.
		winningRoll := middleware.Chance(1.0)
		losingRoll := middleware.Chance(0.0)

		auditorWith := func(roll middleware.MessageAuditor) middleware.MessageAuditor {
			return middleware.AnyOf(
				middleware.NotBot,
				middleware.AllOf(middleware.AuthorID("22222"), roll),
				middleware.ContentContains("bingo"),
			)
		}

		// Human — passes regardless
		Expect(auditorWith(winningRoll).Audit(s, build(authorID("u"), content("hi")))).To(BeTrue())
		// Bot with "bingo" — passes via content branch
		Expect(auditorWith(losingRoll).Audit(s, build(isBot(), authorID("999"), content("bingo")))).To(BeTrue())
		// Bot 22222 on a winning roll — passes via the chance branch
		Expect(auditorWith(winningRoll).Audit(s, build(isBot(), authorID("22222"), content("hi")))).To(BeTrue())
		// Bot 22222 on a losing roll, no "bingo" — blocked
		Expect(auditorWith(losingRoll).Audit(s, build(isBot(), authorID("22222"), content("hi")))).To(BeFalse())
		// Bot 22222 on a losing roll, but has "bingo" — passes via content
		Expect(auditorWith(losingRoll).Audit(s, build(isBot(), authorID("22222"), content("bingo")))).To(BeTrue())
		// Unrelated bot, losing roll, no "bingo" — blocked
		Expect(auditorWith(losingRoll).Audit(s, build(isBot(), authorID("999"), content("hi")))).To(BeFalse())
	})
})
