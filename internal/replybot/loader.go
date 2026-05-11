package replybot

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/andrewgari/starbunk-go/internal/discord"
	"github.com/andrewgari/starbunk-go/internal/middleware"
	"go.yaml.in/yaml/v3"
)

// LoadFile compiles all reply bots defined in a single YAML file.
func LoadFile(path string, messaging discord.MessagingService) ([]*ReplyBot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("replybot: read %s: %w", path, err)
	}

	var cfg FileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("replybot: parse %s: %w", path, err)
	}

	bots := make([]*ReplyBot, 0, len(cfg.ReplyBots))
	for i, botCfg := range cfg.ReplyBots {
		rb, err := compileBot(botCfg, messaging)
		if err != nil {
			return nil, fmt.Errorf("replybot: %s bot[%d] %q: %w", path, i, botCfg.Name, err)
		}
		bots = append(bots, rb)
	}
	return bots, nil
}

// LoadDir compiles all reply bots from every *.yaml / *.yml file in a directory.
// Files are processed in alphabetical order for determinism.
// A parse or compile error in any file is returned immediately.
func LoadDir(dir string, messaging discord.MessagingService) ([]*ReplyBot, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("replybot: read dir %s: %w", dir, err)
	}

	var bots []*ReplyBot
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		ext := filepath.Ext(name)
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		path := filepath.Join(dir, name)
		fileBots, err := LoadFile(path, messaging)
		if err != nil {
			return nil, err
		}
		bots = append(bots, fileBots...)
	}
	return bots, nil
}

// compileBot converts one BotConfig into a live ReplyBot.
func compileBot(cfg BotConfig, messaging discord.MessagingService) (*ReplyBot, error) {
	if cfg.Name == "" {
		return nil, fmt.Errorf("bot name must not be empty")
	}

	globalAuditor, err := buildGlobalAuditor(cfg)
	if err != nil {
		return nil, err
	}

	resolver, err := BuildIdentityResolver(cfg.Identity)
	if err != nil {
		return nil, fmt.Errorf("identity: %w", err)
	}

	var botPool *ResponsePool
	if len(cfg.Responses) > 0 {
		botPool, err = NewResponsePool(cfg.Responses)
		if err != nil {
			return nil, fmt.Errorf("responses: %w", err)
		}
	}

	if len(cfg.Triggers) == 0 {
		log.Printf("replybot: bot %q has no triggers and will never respond", cfg.Name)
	}

	triggers := make([]CompiledTrigger, 0, len(cfg.Triggers))
	for i, tc := range cfg.Triggers {
		ct, err := compileTrigger(tc, i)
		if err != nil {
			return nil, err
		}
		triggers = append(triggers, ct)
	}

	return &ReplyBot{
		Name:             cfg.Name,
		GlobalAuditor:    globalAuditor,
		IdentityResolver: resolver,
		Triggers:         triggers,
		BotPool:          botPool,
		Messaging:        messaging,
	}, nil
}

// compileTrigger converts one TriggerConfig into a CompiledTrigger.
func compileTrigger(tc TriggerConfig, idx int) (CompiledTrigger, error) {
	name := tc.Name
	if name == "" {
		name = fmt.Sprintf("trigger[%d]", idx)
	}

	if tc.Conditions.isEmpty() {
		return CompiledTrigger{}, fmt.Errorf("trigger %q: conditions must not be empty", name)
	}

	condition, err := Compile(tc.Conditions)
	if err != nil {
		return CompiledTrigger{}, fmt.Errorf("trigger %q: %w", name, err)
	}

	var pool *ResponsePool
	if len(tc.Responses) > 0 {
		pool, err = NewResponsePool(tc.Responses)
		if err != nil {
			return CompiledTrigger{}, fmt.Errorf("trigger %q responses: %w", name, err)
		}
	}

	return CompiledTrigger{
		Name:      name,
		Condition: condition,
		Pool:      pool,
	}, nil
}

// buildGlobalAuditor creates the per-bot author filter from ignore_bots / ignore_humans.
func buildGlobalAuditor(cfg BotConfig) (middleware.MessageAuditor, error) {
	switch {
	case cfg.IgnoreBots && cfg.IgnoreHumans:
		log.Printf("replybot: bot %q has ignore_bots=true and ignore_humans=true: it will never respond", cfg.Name)
		return neverAuditor{}, nil
	case cfg.IgnoreBots && !cfg.IgnoreHumans:
		return middleware.NotBot, nil
	case !cfg.IgnoreBots && cfg.IgnoreHumans:
		return middleware.IsBot, nil
	default: // both false
		return alwaysAuditor{}, nil
	}
}
