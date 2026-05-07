package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Run initialized a discord session, registers the provided handlers, and waits for a kill signal.
func Run(botName string, handlers ...any) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN environment variable not set")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("error creating Discord session for %s: %v", botName, err)
	}

	for _, handler := range handlers {
		dg.AddHandler(handler)
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
