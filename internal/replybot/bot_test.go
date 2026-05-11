package replybot_test

import (
	"errors"

	"github.com/andrewgari/starbunk-go/internal/bot"
	"github.com/andrewgari/starbunk-go/internal/middleware"
	"github.com/andrewgari/starbunk-go/internal/replybot"
	"github.com/bwmarrin/discordgo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// mockMessagingService captures SendMessageWithIdentity calls for assertions.
type mockMessagingService struct {
	calls []identityCall
}

type identityCall struct {
	channelID, content, username, avatarURL string
}

func (m *mockMessagingService) SendMessageWithIdentity(channelID, content, username, avatarURL string) (*discordgo.Message, error) {
	m.calls = append(m.calls, identityCall{channelID, content, username, avatarURL})
	return nil, nil
}
func (m *mockMessagingService) SendMessage(string, string) (*discordgo.Message, error) {
	return nil, nil
}
func (m *mockMessagingService) SendComplexMessage(string, *discordgo.MessageSend) (*discordgo.Message, error) {
	return nil, nil
}
func (m *mockMessagingService) ReplyMessage(string, string, string) (*discordgo.Message, error) {
	return nil, nil
}
func (m *mockMessagingService) EditMessage(string, string, string) (*discordgo.Message, error) {
	return nil, nil
}
func (m *mockMessagingService) DeleteMessage(string, string) error { return nil }

// staticIdentity returns an IdentityResolver that always returns the given identity.
func staticIdentity(username, avatarURL string) replybot.IdentityResolver {
	id := bot.Identity{Username: username, AvatarURL: avatarURL}
	return func(_ *discordgo.Session, _ *discordgo.MessageCreate) (bot.Identity, error) {
		return id, nil
	}
}

// errorIdentity returns an IdentityResolver that always errors.
func errorIdentity() replybot.IdentityResolver {
	return func(_ *discordgo.Session, _ *discordgo.MessageCreate) (bot.Identity, error) {
		return bot.Identity{}, errors.New("identity error")
	}
}

// fixedPool creates a ResponsePool that always returns the given string.
func fixedPool(response string) *replybot.ResponsePool {
	p, _ := replybot.NewResponsePoolWithRng([]string{response}, func(_ int) int { return 0 })
	return p
}

// alwaysAuditor is a test auditor that always passes.
type testAlwaysAuditor struct{}

func (testAlwaysAuditor) Audit(_ *discordgo.Session, _ *discordgo.MessageCreate) bool { return true }

// neverAuditor is a test auditor that always fails.
type testNeverAuditor struct{}

func (testNeverAuditor) Audit(_ *discordgo.Session, _ *discordgo.MessageCreate) bool { return false }

var _ = Describe("ReplyBot.Handle", func() {
	var (
		s       *discordgo.Session
		mock    *mockMessagingService
		trigger replybot.CompiledTrigger
	)

	BeforeEach(func() {
		s = testSession()
		mock = &mockMessagingService{}
		trigger = replybot.CompiledTrigger{
			Name:      "test-trigger",
			Condition: testAlwaysAuditor{},
			Pool:      fixedPool("trigger response"),
		}
	})

	It("sends a message when the trigger matches", func() {
		rb := &replybot.ReplyBot{
			Name:             "test-bot",
			GlobalAuditor:    testAlwaysAuditor{},
			IdentityResolver: staticIdentity("BotName", "https://example.com/avatar.png"),
			Triggers:         []replybot.CompiledTrigger{trigger},
			Messaging:        mock,
		}
		rb.Handle(s, buildMsg(withContent("hello"), withAuthorID("u1")))
		Expect(mock.calls).To(HaveLen(1))
		Expect(mock.calls[0].content).To(Equal("trigger response"))
		Expect(mock.calls[0].username).To(Equal("BotName"))
	})

	It("does not send when the global auditor blocks", func() {
		rb := &replybot.ReplyBot{
			Name:             "test-bot",
			GlobalAuditor:    testNeverAuditor{},
			IdentityResolver: staticIdentity("BotName", "https://example.com/avatar.png"),
			Triggers:         []replybot.CompiledTrigger{trigger},
			Messaging:        mock,
		}
		rb.Handle(s, buildMsg(withContent("hello")))
		Expect(mock.calls).To(BeEmpty())
	})

	It("drops bot messages when GlobalAuditor is NotBot", func() {
		rb := &replybot.ReplyBot{
			Name:             "test-bot",
			GlobalAuditor:    middleware.NotBot,
			IdentityResolver: staticIdentity("BotName", "https://example.com/avatar.png"),
			Triggers:         []replybot.CompiledTrigger{trigger},
			Messaging:        mock,
		}
		rb.Handle(s, buildMsg(withBot()))
		Expect(mock.calls).To(BeEmpty())
	})

	It("fires the first matching trigger and stops", func() {
		trigger2 := replybot.CompiledTrigger{
			Name:      "second",
			Condition: testAlwaysAuditor{},
			Pool:      fixedPool("second response"),
		}
		rb := &replybot.ReplyBot{
			Name:             "test-bot",
			GlobalAuditor:    testAlwaysAuditor{},
			IdentityResolver: staticIdentity("BotName", "https://example.com/avatar.png"),
			Triggers:         []replybot.CompiledTrigger{trigger, trigger2},
			Messaging:        mock,
		}
		rb.Handle(s, buildMsg(withContent("hello")))
		Expect(mock.calls).To(HaveLen(1))
		Expect(mock.calls[0].content).To(Equal("trigger response"))
	})

	It("falls back to bot pool when trigger has no pool", func() {
		triggerNoPool := replybot.CompiledTrigger{
			Name:      "no-pool",
			Condition: testAlwaysAuditor{},
			Pool:      nil,
		}
		rb := &replybot.ReplyBot{
			Name:             "test-bot",
			GlobalAuditor:    testAlwaysAuditor{},
			IdentityResolver: staticIdentity("BotName", "https://example.com/avatar.png"),
			Triggers:         []replybot.CompiledTrigger{triggerNoPool},
			BotPool:          fixedPool("bot pool response"),
			Messaging:        mock,
		}
		rb.Handle(s, buildMsg(withContent("hello")))
		Expect(mock.calls).To(HaveLen(1))
		Expect(mock.calls[0].content).To(Equal("bot pool response"))
	})

	It("sends nothing when no trigger matches", func() {
		triggerNever := replybot.CompiledTrigger{
			Name:      "never",
			Condition: testNeverAuditor{},
			Pool:      fixedPool("never"),
		}
		rb := &replybot.ReplyBot{
			Name:             "test-bot",
			GlobalAuditor:    testAlwaysAuditor{},
			IdentityResolver: staticIdentity("BotName", "https://example.com/avatar.png"),
			Triggers:         []replybot.CompiledTrigger{triggerNever},
			Messaging:        mock,
		}
		rb.Handle(s, buildMsg(withContent("hello")))
		Expect(mock.calls).To(BeEmpty())
	})

	It("sends nothing when trigger and bot both lack a pool", func() {
		triggerNoPool := replybot.CompiledTrigger{
			Name:      "no-pool",
			Condition: testAlwaysAuditor{},
			Pool:      nil,
		}
		rb := &replybot.ReplyBot{
			Name:             "test-bot",
			GlobalAuditor:    testAlwaysAuditor{},
			IdentityResolver: staticIdentity("BotName", "https://example.com/avatar.png"),
			Triggers:         []replybot.CompiledTrigger{triggerNoPool},
			BotPool:          nil,
			Messaging:        mock,
		}
		rb.Handle(s, buildMsg(withContent("hello")))
		Expect(mock.calls).To(BeEmpty())
	})

	It("uses fallback identity when resolver errors", func() {
		rb := &replybot.ReplyBot{
			Name:             "test-bot",
			GlobalAuditor:    testAlwaysAuditor{},
			IdentityResolver: errorIdentity(),
			Triggers:         []replybot.CompiledTrigger{trigger},
			Messaging:        mock,
		}
		// Should not panic; sends with empty/fallback identity.
		Expect(func() {
			rb.Handle(s, buildMsg(withContent("hello")))
		}).NotTo(Panic())
		// Message was still attempted (using fallback identity).
		Expect(mock.calls).To(HaveLen(1))
	})
})
