package bot

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/andrewgari/starbunk-go/internal/middleware"
	"github.com/bwmarrin/discordgo"
)

// Run initialises a Discord session, registers the provided handlers, and waits
// for a kill signal.
//
// auditor is applied to every MessageCreate handler before it is invoked. Any
// handler whose signature is func(*discordgo.Session, *discordgo.MessageCreate)
// is automatically wrapped — no message reaches it without passing audit. Other
// event types (voice state updates, reactions, etc.) are registered directly and
// are not subject to message audit.
func Run(botName string, auditor middleware.MessageAuditor, handlers ...any) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		slog.Error("DISCORD_TOKEN environment variable not set", "bot", botName)
		os.Exit(1)
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		slog.Error("error creating Discord session", "bot", botName, "err", err)
		os.Exit(1)
	}

	for _, h := range handlers {
		if msgHandler, ok := h.(func(*discordgo.Session, *discordgo.MessageCreate)); ok {
			dg.AddHandler(auditedHandler(auditor, msgHandler))
		} else {
			dg.AddHandler(h)
		}
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	err = dg.Open()
	if err != nil {
		slog.Error("error opening connection", "bot", botName, "err", err)
		os.Exit(1)
	}

	slog.Info("bot running, press CTRL-C to exit", "bot", botName)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if err := dg.Close(); err != nil {
		slog.Error("error closing Discord session", "bot", botName, "err", err)
	}
}

func auditedHandler(a middleware.MessageAuditor, h func(*discordgo.Session, *discordgo.MessageCreate)) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if a.Audit(s, m) {
			h(s, m)
		}
	}
}
