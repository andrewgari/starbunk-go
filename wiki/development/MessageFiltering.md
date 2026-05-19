# Message Audit Abstraction

> **Status:** Implemented — `internal/middleware`
> **Relates to:** [[../infrastructure/Architecture|Architecture]], [[../bots/BunkBot|BunkBot]], [[../bots/BlueBot|BlueBot]], [[../bots/CovaBot|CovaBot]]

---

## Problem

Every bot currently has one inline guard copied verbatim:

```go
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == s.State.User.ID {
        return
    }
    // bot logic
}
```

This is insufficient for two reasons:

1. **It is opt-in.** A bot author can forget to write it. Nothing in the framework enforces it.
2. **It does not scale.** Each bot will need a different set of acceptance rules (self-exclusion, bot exclusion, guild-only, content requirements, BunkBot-specific rule chains). Without a shared abstraction, every bot reimplements the same logic differently, making it untestable in isolation and impossible to compose.

The fix is not a helper function bots can choose to call. The fix is **making audit mandatory at the framework level**, so no message can reach any handler without passing through the bot's declared policy.

---

## Design Principle

> Every `MessageCreate` event passes through the bot's `MessageAuditor` before any handler sees it. This is enforced by `bot.Run`, not by bot authors.

Bots declare **what they accept** (their auditor). The framework handles **ensuring it runs** on every message. These two concerns are separated.

---

## Core Interface: `MessageAuditor`

```go
// internal/middleware/auditor.go

// MessageAuditor is the mandatory evaluation gate for incoming Discord messages.
// bot.Run requires one. The framework automatically applies it to every
// MessageCreate event before invoking any registered handler.
//
// Audit returns true if the message should be processed, false if it should
// be dropped silently.
type MessageAuditor interface {
    Audit(s *discordgo.Session, m *discordgo.MessageCreate) bool
}
```

---

## Framework Enforcement: `bot.Run` Signature Change

`bot.Run` gains a required `MessageAuditor` parameter. There is no way to start a bot without supplying one.

```go
// internal/bot/bot.go

// Run initialises the bot with a mandatory message auditor.
//
// Every handler whose signature is func(*discordgo.Session, *discordgo.MessageCreate)
// is automatically wrapped: the auditor runs first, and the handler is only
// called if auditor.Audit returns true.
//
// Handlers for other event types (voice state updates, reactions, etc.) are
// registered directly and are not subject to message audit.
func Run(name string, auditor middleware.MessageAuditor, handlers ...any)
```

### How wrapping works inside `bot.Run`

```go
for _, h := range handlers {
    switch typed := h.(type) {
    case func(*discordgo.Session, *discordgo.MessageCreate):
        // Enforce the audit for every message handler.
        dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
            if auditor.Audit(s, m) {
                typed(s, m)
            }
        })
    default:
        // Voice, reactions, guild events — registered without audit.
        dg.AddHandler(h)
    }
}
```

Bot authors cannot bypass this. Passing a raw `func(*discordgo.Session, *discordgo.MessageCreate)` to `bot.Run` always goes through the auditor. There is no `Filtered()` helper to forget to call.

---

## Primitives

Atomic building blocks in `internal/middleware`. Each is a package-level variable of a small unexported struct that satisfies `MessageAuditor`. They are stateless and allocation-free.

| Var | Drops when… |
|---|---|
| `NotSelf` | `m.Author.ID == s.State.User.ID` |
| `NotBot` | `m.Author.Bot == true` |
| `HasContent` | `strings.TrimSpace(m.Content) == ""` |
| `GuildOnly` | `m.GuildID == ""` (i.e. a DM) |
| `DMOnly` | `m.GuildID != ""` (i.e. a guild message) |

```go
var (
    NotSelf    MessageAuditor = notSelfAuditor{}
    NotBot     MessageAuditor = notBotAuditor{}
    HasContent MessageAuditor = hasContentAuditor{}
    GuildOnly  MessageAuditor = guildOnlyAuditor{}
    DMOnly     MessageAuditor = dmOnlyAuditor{}
)
```

---

## Combinators

Build composite auditors from primitives. All return `MessageAuditor`.

```go
// AllOf passes only when every child auditor passes. Short-circuits on first failure.
func AllOf(auditors ...MessageAuditor) MessageAuditor

// AnyOf passes when at least one child auditor passes. Short-circuits on first success.
func AnyOf(auditors ...MessageAuditor) MessageAuditor

// Not inverts an auditor.
func Not(a MessageAuditor) MessageAuditor
```

---

## Per-Bot Auditors

Each bot declares its auditor inline in `main.go`. The auditor is the first thing visible when reading the bot's entry point — it is the bot's policy statement.

### BlueBot

Never triggers off itself or any other bot. Guild messages only. Must have content.

```go
// cmd/bluebot/main.go
func main() {
    bot.Run("BlueBot",
        middleware.AllOf(
            middleware.NotSelf,
            middleware.NotBot,
            middleware.GuildOnly,
            middleware.HasContent,
        ),
        messageCreate,
    )
}
```

### CovaBot

Same base policy as BlueBot. Conversational — should never respond to itself or other bots.

```go
// cmd/covabot/main.go
func main() {
    bot.Run("CovaBot",
        middleware.AllOf(
            middleware.NotSelf,
            middleware.NotBot,
            middleware.GuildOnly,
            middleware.HasContent,
        ),
        messageCreate,
    )
}
```

### BunkBot

More permissive. Does not exclude bot messages (needed for moderation and persona replies). May support DMs. Rule set will grow as features are added — start minimal and extend.

```go
// cmd/bunkbot/main.go
func main() {
    bot.Run("BunkBot",
        middleware.AllOf(
            middleware.NotSelf,   // never respond to itself
            middleware.HasContent, // ignore empty/system messages
            // additional rules added here as BunkBot features are implemented:
            // e.g. middleware.NotUser("some-id"), cooldowns, channel allow-lists
        ),
        messageCreate,
    )
}
```

### DJCova / RatBot

These bots also require an auditor. As features are implemented, the appropriate policy will be chosen. Until then, a minimal auditor:

```go
middleware.AllOf(middleware.NotSelf, middleware.HasContent)
```

---

## What Happens to the Inline Guard

The existing `if m.Author.ID == s.State.User.ID { return }` guard in every `messageCreate` function is **deleted**. The auditor replaces it. After this change, `messageCreate` functions contain only business logic — they assume the message has already been admitted.

---

## Extending the System

### New stateless primitive

```go
// Drop messages from a specific user ID.
type notUserAuditor struct{ userID string }

func (a notUserAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
    return m.Author.ID != a.userID
}

// Constructor, not a var, because it takes a parameter.
func NotUser(userID string) MessageAuditor {
    return notUserAuditor{userID: userID}
}
```

### New stateful auditor (e.g. per-user cooldown)

Stateful auditors are constructed with `New<Name>(...)` and use a mutex-protected store internally. They satisfy `MessageAuditor` via a pointer receiver.

```go
type cooldownAuditor struct {
    mu       sync.Mutex
    lastSeen map[string]time.Time
    window   time.Duration
}

func NewCooldown(window time.Duration) MessageAuditor {
    return &cooldownAuditor{
        lastSeen: make(map[string]time.Time),
        window:   window,
    }
}

func (a *cooldownAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
    a.mu.Lock()
    defer a.mu.Unlock()
    key := m.Author.ID + ":" + m.ChannelID
    if t, ok := a.lastSeen[key]; ok && time.Since(t) < a.window {
        return false
    }
    a.lastSeen[key] = time.Now()
    return true
}
```

To use: `middleware.AllOf(middleware.NotSelf, middleware.NewCooldown(5*time.Second))`.

---

## Testing Strategy

Primitives are pure functions over `discordgo.MessageCreate`. Tests construct minimal structs — no HTTP, no WebSocket, no Discord mock needed.

```go
// internal/middleware/auditor_test.go

Describe("NotSelf", func() {
    var s *discordgo.Session

    BeforeEach(func() {
        s = &discordgo.Session{State: &discordgo.State{}}
        s.State.User = &discordgo.User{ID: "bot-id"}
    })

    It("passes messages from other users", func() {
        m := messageFrom("user-id")
        Expect(middleware.NotSelf.Audit(s, m)).To(BeTrue())
    })

    It("drops messages from the bot itself", func() {
        m := messageFrom("bot-id")
        Expect(middleware.NotSelf.Audit(s, m)).To(BeFalse())
    })
})

Describe("AllOf", func() {
    It("passes when all children pass", func() {
        auditor := middleware.AllOf(middleware.NotSelf, middleware.HasContent)
        s := sessionWithID("bot-id")
        m := messageFromWithContent("user-id", "hello")
        Expect(auditor.Audit(s, m)).To(BeTrue())
    })

    It("drops when any child fails", func() {
        auditor := middleware.AllOf(middleware.NotSelf, middleware.HasContent)
        s := sessionWithID("bot-id")
        m := messageFromWithContent("user-id", "   ")
        Expect(auditor.Audit(s, m)).To(BeFalse())
    })
})
```

---

## File Layout

```
internal/
  middleware/
    auditor.go        # MessageAuditor interface, primitives, combinators
    auditor_test.go   # Ginkgo suite covering all primitives and combinators
  bot/
    bot.go            # Updated: Run() takes MessageAuditor, auto-wraps handlers
```

Auditor compositions (the per-bot policy) live in each `cmd/<bot>/main.go`, not in `internal/`, because they express bot-specific business rules.

---

## Implementation Checklist

- [x] Create `internal/middleware/auditor.go` — `MessageAuditor` interface, `AllOf`, `AnyOf`, `Not`
- [x] `author.go` — `NotSelf`, `NotBot`, `IsBot`, `AuthorID`, `NotAuthorID`, `AuthorNamed`, `AuthorHasRole`
- [x] `content.go` — `HasContent`, `ContentContains`, `ContentMatches`, `HasAttachment`
- [x] `context.go` — `GuildOnly`, `DMOnly`, `InChannel`, `OnWeekdays`
- [x] `random.go` — `Chance`
- [x] Write Ginkgo suite `internal/middleware/auditor_test.go` (48 specs)
- [x] Update `internal/bot/bot.go` — `Run` requires `MessageAuditor`; auto-wraps `MessageCreate` handlers
- [x] Update each bot's `main.go` to declare an auditor
  - [x] `cmd/bluebot` — `AllOf(NotSelf, NotBot, GuildOnly, HasContent)`
  - [x] `cmd/covabot` — `AllOf(NotSelf, NotBot, GuildOnly, HasContent)`
  - [x] `cmd/bunkbot` — `AllOf(NotSelf, HasContent)`
  - [x] `cmd/djcova` — `AllOf(NotSelf, HasContent)`
  - [x] `cmd/ratbot` — `AllOf(NotSelf, HasContent)`
- [x] Delete inline self-guards from all `messageCreate` functions
- [x] Update [[../infrastructure/Architecture|Architecture]] wiki page
- [x] Add entry to [[../Changelog|Changelog]]

---

---

## Tier 2: Strategy-Level Conditions (`internal/replybot`)

For bots that use the `replybot.Bot` dispatcher (e.g. BlueBot, and eventually BunkBot), a second
tier of filtering exists at the **individual strategy level**. This is necessary when different
strategies within the same bot need conflicting filter policies — for example, a "BotBot" strategy
that responds only to bots, alongside human-only strategies.

### How it works

`Bot.Handle()` checks whether a strategy implements the optional `ConditionedStrategy` interface
before calling `ShouldTrigger`. If the condition fails, the message is silently skipped for that
strategy and evaluation continues with the next one.

```go
// internal/replybot/strategy.go

// ConditionedStrategy is an optional extension of Strategy.
type ConditionedStrategy interface {
    Strategy
    Condition() middleware.MessageAuditor
}
```

`Bot.Handle` accepts a `*discordgo.Session` so conditions like `AuthorHasRole` can inspect guild
state:

```go
// internal/replybot/bot.go
func (b *Bot) Handle(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate)
```

### WithCondition — compose at the call site

`WithCondition` wraps any existing `Strategy` with a condition without requiring changes to the
strategy struct:

```go
// internal/replybot/condition.go
func WithCondition(cond middleware.MessageAuditor, s Strategy) Strategy
```

### Example — BunkBot with mixed strategies

```go
replybot.NewBot(sender,
    // BotBot: responds only when the author is a bot
    replybot.WithCondition(middleware.IsBot, botBotStrategy),

    // Human-only strategy: responds only to non-bot authors
    replybot.WithCondition(middleware.NotBot, someHumanStrategy),

    // Unconditioned: no extra filtering beyond the bot-level auditor
    anotherStrategy,
)
```

### Tier summary

| Tier | Mechanism | Where declared | Example |
|------|-----------|----------------|---------|
| 1 | `MessageAuditor` in `bot.Run()` | `cmd/<bot>/main.go` | `NotSelf`, `NotBot` (BlueBot) |
| 2 | `ConditionedStrategy` / `WithCondition` | Strategy construction | `IsBot` for BotBot |

---

## See Also

- [[../infrastructure/Architecture|Architecture]] — shared library overview and bot pattern
- [[../bots/BunkBot|BunkBot]] — most complex auditor policy
- [[../bots/BlueBot|BlueBot]] — standard strict auditor policy
- [[../bots/CovaBot|CovaBot]] — standard strict auditor policy
