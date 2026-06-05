package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/andrewgari/starbunk-go/internal/bot"
	"github.com/andrewgari/starbunk-go/internal/discord"
	"github.com/andrewgari/starbunk-go/internal/llm"
	"github.com/andrewgari/starbunk-go/internal/memory"
	"github.com/andrewgari/starbunk-go/internal/middleware"
	"github.com/bwmarrin/discordgo"
)

var auditor = middleware.AllOf(
	middleware.NotSelf,
	middleware.NotBot,
	middleware.GuildOnly,
	middleware.HasContent,
)

type Handler struct {
	llms   llm.Registry
	memory memory.Service
}

func main() {
	cfg := llm.ConfigFromEnv()
	llms, err := llm.NewRegistry(cfg)
	if err != nil {
		slog.Error("failed to init llm registry", "err", err)
		os.Exit(1)
	}

	dbConnStr := os.Getenv("POSTGRES_CONN_STR")
	if dbConnStr == "" {
		//nolint:gosec
		dbConnStr = "postgres://starbunk:starbunk@starbunk-go-postgres:5432/starbunk_memory?sslmode=disable"
	}

	store, err := memory.NewPGStore(dbConnStr)
	if err != nil {
		slog.Error("failed to init memory store", "err", err)
		os.Exit(1)
	}
	defer func() { _ = store.Close() }()

	memService := memory.NewService(store, llms)

	h := &Handler{
		llms:   llms,
		memory: memService,
	}

	bot.Run("CovaBot", auditor, h.messageCreate)
}

func (h *Handler) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "ping covabot" {
		sender := discord.NewMessageService(s)
		_, err := sender.SendMessage(m.ChannelID, "Pong from covabot!")
		if err != nil {
			slog.Error("failed to send message", "bot", "covabot", "err", err)
		}
		return
	}

	ctx := context.Background()

	// 1. Asynchronously extract and save memory
	h.memory.ExtractAndSave(ctx, m.Author.ID, m.Content)

	// 2. Recall context
	memContext, err := h.memory.Recall(ctx, m.Content)
	if err != nil {
		slog.Warn("failed to recall memory", "err", err)
	}

	// 3. Generate response
	highLLM := h.llms.High()
	if highLLM == nil {
		slog.Warn("no high LLM available to respond")
		return
	}

	systemPrompt := "You are CovaBot, a helpful AI personality. Respond to the user conversationally."
	if memContext != "" {
		systemPrompt += fmt.Sprintf("\n\nRelevant past memories/facts:\n%s", memContext)
	}

	req := llm.GenerateRequest{
		Messages: []llm.Message{
			{Role: llm.RoleSystem, Content: systemPrompt},
			{Role: llm.RoleUser, Content: m.Content},
		},
	}

	resp, err := highLLM.Generate(ctx, req)
	if err != nil {
		slog.Error("failed to generate response", "err", err)
		return
	}

	sender := discord.NewMessageService(s)
	_, err = sender.SendMessage(m.ChannelID, resp.Text)
	if err != nil {
		slog.Error("failed to send message", "err", err)
	}
}
