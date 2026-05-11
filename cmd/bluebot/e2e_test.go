package main

import (
	"github.com/andrewgari/starbunk-go/internal/testenv"
	"github.com/bwmarrin/discordgo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BlueBot e2e", func() {
	// Auditor policy: AllOf(NotSelf, NotBot, GuildOnly, HasContent)

	var harness *testenv.BotHarness

	BeforeEach(func() {
		harness = testenv.NewBotHarness("bluebot-bot-id", auditor, messageCreate)
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
				testenv.WithAuthorID("bluebot-bot-id"), // matches harness botID
				testenv.WithGuildID(testenv.DefaultGuildID),
				testenv.WithContent("hello"),
			),
			wantAudit: false,
		}),

		Entry("bot author is blocked by NotBot", scenario{
			desc: "bot author blocked by NotBot",
			msg: testenv.NewMessage(
				testenv.WithAuthorID("other-bot-id"),
				testenv.WithBot(),
				testenv.WithGuildID(testenv.DefaultGuildID),
				testenv.WithContent("hello"),
			),
			wantAudit: false,
		}),

		Entry("DM is blocked by GuildOnly", scenario{
			desc: "DM blocked by GuildOnly",
			msg: testenv.NewMessage(
				testenv.WithAuthorID("user-123"),
				testenv.WithContent("hello"),
				// no WithGuildID → GuildID empty → GuildOnly fails
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

		// — handler round-trip scenarios —

		Entry("ping bluebot triggers pong", scenario{
			desc: "ping bluebot → Pong from bluebot!",
			msg: testenv.NewMessage(
				testenv.WithAuthorID("user-123"),
				testenv.WithGuildID(testenv.DefaultGuildID),
				testenv.WithContent("ping bluebot"),
			),
			wantAudit:   true,
			wantMsgs:    1,
			wantContent: "Pong from bluebot!",
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
