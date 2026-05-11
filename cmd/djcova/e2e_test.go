package main

import (
	"github.com/andrewgari/starbunk-go/internal/testenv"
	"github.com/bwmarrin/discordgo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DJCova e2e", func() {
	// Auditor policy: AllOf(NotSelf, HasContent)
	// Notably: bot authors and DMs are ALLOWED (no NotBot, no GuildOnly).

	var harness *testenv.BotHarness

	BeforeEach(func() {
		harness = testenv.NewBotHarness("djcova-bot-id", auditor, messageCreate)
	})

	type scenario struct {
		desc        string
		msg         *discordgo.MessageCreate
		wantAudit   bool
		wantMsgs    int
		wantContent string
	}

	DescribeTable("message routing",
		func(sc scenario) {
			result := harness.Run(sc.msg)

			Expect(result.AuditPassed).To(Equal(sc.wantAudit),
				testenv.FormatScenarioFailure(sc.desc, result))

			if sc.wantAudit {
				Expect(result.Messages).To(HaveLen(sc.wantMsgs),
					testenv.FormatScenarioFailure(sc.desc, result))
				if sc.wantContent != "" && sc.wantMsgs > 0 {
					Expect(result.Messages[0].Content).To(Equal(sc.wantContent),
						testenv.FormatScenarioFailure(sc.desc, result))
				}
			}
		},

		// — auditor gate scenarios —

		Entry("human guild message passes audit", scenario{
			desc: "human guild message passes audit",
			msg: testenv.NewMessage(
				testenv.WithAuthorID("user-123"),
				testenv.WithGuildID(testenv.DefaultGuildID),
				testenv.WithContent("hello"),
			),
			wantAudit: true,
		}),

		Entry("self message is blocked by NotSelf", scenario{
			desc: "self message blocked by NotSelf",
			msg: testenv.NewMessage(
				testenv.WithAuthorID("djcova-bot-id"), // matches harness botID
				testenv.WithGuildID(testenv.DefaultGuildID),
				testenv.WithContent("hello"),
			),
			wantAudit: false,
		}),

		Entry("empty content is blocked by HasContent", scenario{
			desc: "empty content blocked by HasContent",
			msg: testenv.NewMessage(
				testenv.WithAuthorID("user-123"),
				testenv.WithGuildID(testenv.DefaultGuildID),
				testenv.WithContent(""),
			),
			wantAudit: false,
		}),

		// DJCova intentionally allows bot authors (no NotBot check).
		Entry("bot author passes audit (no NotBot policy)", scenario{
			desc: "bot author passes (DJCova allows bots)",
			msg: testenv.NewMessage(
				testenv.WithAuthorID("other-bot-id"),
				testenv.WithBot(),
				testenv.WithGuildID(testenv.DefaultGuildID),
				testenv.WithContent("hello"),
			),
			wantAudit: true,
		}),

		// DJCova intentionally allows DMs (no GuildOnly check).
		Entry("DM passes audit (no GuildOnly policy)", scenario{
			desc: "DM passes (DJCova allows DMs)",
			msg: testenv.NewMessage(
				testenv.WithAuthorID("user-123"),
				testenv.WithContent("hello"),
				// no WithGuildID → DM context
			),
			wantAudit: true,
		}),

		// — handler round-trip scenarios —

		Entry("ping djcova triggers pong", scenario{
			desc: "ping djcova → Pong from djcova!",
			msg: testenv.NewMessage(
				testenv.WithAuthorID("user-123"),
				testenv.WithGuildID(testenv.DefaultGuildID),
				testenv.WithContent("ping djcova"),
			),
			wantAudit:   true,
			wantMsgs:    1,
			wantContent: "Pong from djcova!",
		}),

		Entry("unrecognised content passes audit but sends no reply", scenario{
			desc: "unrecognised content → no reply",
			msg: testenv.NewMessage(
				testenv.WithAuthorID("user-123"),
				testenv.WithGuildID(testenv.DefaultGuildID),
				testenv.WithContent("hello world"),
			),
			wantAudit: true,
			wantMsgs:  0,
		}),
	)
})
