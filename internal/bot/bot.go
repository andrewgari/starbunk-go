package bot

import (
	"fmt"
	"log"
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
		log.Fatal("DISCORD_TOKEN environment variable not set")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("error creating Discord session for %s: %v", botName, err)
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
		log.Fatalf("error opening connection for %s: %v", botName, err)
	}

	fmt.Printf("%s is now running. Press CTRL-C to exit.\n", botName)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if err := dg.Close(); err != nil {
		log.Printf("error closing Discord session for %s: %v", botName, err)
	}
}

func auditedHandler(a middleware.MessageAuditor, h func(*discordgo.Session, *discordgo.MessageCreate)) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if a.Audit(s, m) {
			h(s, m)
		}
	}
}
