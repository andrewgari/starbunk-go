package replybot_test

import (
	"context"
	"testing"

	"github.com/andrewgari/starbunk-go/internal/discord"
	"github.com/andrewgari/starbunk-go/internal/middleware"
	"github.com/andrewgari/starbunk-go/internal/replybot"
	"github.com/bwmarrin/discordgo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestReplybot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Replybot Suite")
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
			Author: &discordgo.User{},
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
func isBot() msgOpt {
	return func(m *discordgo.MessageCreate) { m.Author.Bot = true }
}
func withContent(c string) msgOpt {
	return func(m *discordgo.MessageCreate) { m.Content = c }
}
func inGuild(id string) msgOpt {
	return func(m *discordgo.MessageCreate) { m.GuildID = id }
}

// stubSender captures the channel and content of the most recent SendMessage call.
type stubSender struct {
	lastChannel string
	lastContent string
	callCount   int
}

func (s *stubSender) SendMessage(channelID, content string) (*discordgo.Message, error) {
	s.lastChannel = channelID
	s.lastContent = content
	s.callCount++
	return &discordgo.Message{}, nil
}
func (s *stubSender) SendComplexMessage(_ string, _ *discordgo.MessageSend) (*discordgo.Message, error) {
	return nil, nil
}
func (s *stubSender) ReplyMessage(_, _, _ string) (*discordgo.Message, error) { return nil, nil }
func (s *stubSender) SendMessageWithIdentity(_, _, _, _ string) (*discordgo.Message, error) {
	return nil, nil
}
func (s *stubSender) EditMessage(_, _, _ string) (*discordgo.Message, error) { return nil, nil }
func (s *stubSender) DeleteMessage(_, _ string) error                        { return nil }

// Verify stubSender satisfies MessagingService at compile time.
var _ discord.MessagingService = (*stubSender)(nil)

// stubStrategy records calls and returns a configurable trigger result.
type stubStrategy struct {
	name          string
	triggerResult bool
	response      string
	triggerCalls  int
}

func (s *stubStrategy) Name() string { return s.name }
func (s *stubStrategy) ShouldTrigger(_ context.Context, _ *discordgo.MessageCreate) bool {
	s.triggerCalls++
	return s.triggerResult
}
func (s *stubStrategy) Response(_ context.Context, _ *discordgo.MessageCreate) string {
	return s.response
}

// — specs —

var _ = Describe("Bot.Handle", func() {
	var (
		sess   *discordgo.Session
		sender *stubSender
	)

	BeforeEach(func() {
		sess = session("bot-id")
		sender = &stubSender{}
	})

	Context("unconditioned strategies", func() {
		It("calls ShouldTrigger on an unconditioned strategy for every message", func() {
			strat := &stubStrategy{name: "s", triggerResult: false}
			bot := replybot.NewBot(sender, strat)

			bot.Handle(context.Background(), sess, build(authorID("u"), inGuild("g")))
			Expect(strat.triggerCalls).To(Equal(1))
		})

		It("sends the response of the first triggering strategy", func() {
			strat := &stubStrategy{name: "s", triggerResult: true, response: "pong"}
			bot := replybot.NewBot(sender, strat)

			bot.Handle(context.Background(), sess, build(authorID("u"), withContent("ping"), inGuild("g")))
			Expect(sender.lastContent).To(Equal("pong"))
		})

		It("stops after the first match (first-match-wins)", func() {
			first := &stubStrategy{name: "first", triggerResult: true, response: "first"}
			second := &stubStrategy{name: "second", triggerResult: true, response: "second"}
			bot := replybot.NewBot(sender, first, second)

			bot.Handle(context.Background(), sess, build(authorID("u")))
			Expect(sender.lastContent).To(Equal("first"))
			Expect(second.triggerCalls).To(Equal(0))
		})
	})

	Context("ConditionedStrategy via WithCondition", func() {
		It("skips ShouldTrigger when the condition fails", func() {
			strat := &stubStrategy{name: "bot-only", triggerResult: true}
			bot := replybot.NewBot(sender, replybot.WithCondition(middleware.IsBot, strat))

			// Human message — IsBot condition fails → ShouldTrigger must not be called
			bot.Handle(context.Background(), sess, build(authorID("u")))
			Expect(strat.triggerCalls).To(Equal(0))
			Expect(sender.callCount).To(Equal(0))
		})

		It("calls ShouldTrigger when the condition passes", func() {
			strat := &stubStrategy{name: "bot-only", triggerResult: true, response: "hello bot"}
			bot := replybot.NewBot(sender, replybot.WithCondition(middleware.IsBot, strat))

			// Bot message — IsBot condition passes
			bot.Handle(context.Background(), sess, build(authorID("other-bot"), isBot()))
			Expect(strat.triggerCalls).To(Equal(1))
			Expect(sender.lastContent).To(Equal("hello bot"))
		})

		It("falls through to the next strategy when a condition fails", func() {
			botOnly := &stubStrategy{name: "bot-only", triggerResult: true, response: "bot reply"}
			humanOnly := &stubStrategy{name: "human-only", triggerResult: true, response: "human reply"}
			bot := replybot.NewBot(sender,
				replybot.WithCondition(middleware.IsBot, botOnly),
				replybot.WithCondition(middleware.NotBot, humanOnly),
			)

			// Human message: botOnly condition fails, humanOnly condition passes
			bot.Handle(context.Background(), sess, build(authorID("u")))
			Expect(botOnly.triggerCalls).To(Equal(0))
			Expect(humanOnly.triggerCalls).To(Equal(1))
			Expect(sender.lastContent).To(Equal("human reply"))

			// Reset
			botOnly.triggerCalls = 0
			humanOnly.triggerCalls = 0
			sender.lastContent = ""

			// Bot message: botOnly condition passes, humanOnly is never reached
			bot.Handle(context.Background(), sess, build(authorID("other-bot"), isBot()))
			Expect(botOnly.triggerCalls).To(Equal(1))
			Expect(humanOnly.triggerCalls).To(Equal(0))
			Expect(sender.lastContent).To(Equal("bot reply"))
		})

		It("mixes conditioned and unconditioned strategies correctly", func() {
			conditioned := &stubStrategy{name: "bot-only", triggerResult: false}
			fallback := &stubStrategy{name: "fallback", triggerResult: true, response: "fallback"}
			bot := replybot.NewBot(sender,
				replybot.WithCondition(middleware.IsBot, conditioned),
				fallback,
			)

			// Human: conditioned is skipped (condition fails), fallback is reached
			bot.Handle(context.Background(), sess, build(authorID("u")))
			Expect(conditioned.triggerCalls).To(Equal(0))
			Expect(fallback.triggerCalls).To(Equal(1))
		})

		It("uses AllOf conditions for multi-criteria filtering", func() {
			strat := &stubStrategy{name: "bot-guild", triggerResult: true, response: "ok"}
			cond := middleware.AllOf(middleware.IsBot, middleware.GuildOnly)
			bot := replybot.NewBot(sender, replybot.WithCondition(cond, strat))

			// Bot in DM → condition fails (not GuildOnly)
			bot.Handle(context.Background(), sess, build(isBot()))
			Expect(strat.triggerCalls).To(Equal(0))

			// Bot in guild → condition passes
			bot.Handle(context.Background(), sess, build(isBot(), inGuild("g")))
			Expect(strat.triggerCalls).To(Equal(1))
			Expect(sender.lastContent).To(Equal("ok"))
		})
	})

	Context("NotSelf via bot-level auditor (tier-1 integration)", func() {
		It("NotSelf condition skips messages from the bot itself", func() {
			strat := &stubStrategy{name: "self-check", triggerResult: true, response: "oops"}
			bot := replybot.NewBot(sender, replybot.WithCondition(middleware.NotSelf, strat))

			// Self message → condition fails
			bot.Handle(context.Background(), sess, build(authorID("bot-id")))
			Expect(strat.triggerCalls).To(Equal(0))

			// Other user → condition passes
			bot.Handle(context.Background(), sess, build(authorID("other-user")))
			Expect(strat.triggerCalls).To(Equal(1))
		})
	})
})
