package replybot

import (
	"fmt"
	"math/rand/v2"

	"github.com/andrewgari/starbunk-go/internal/bot"
	"github.com/bwmarrin/discordgo"
)

// IdentityResolver resolves the bot persona to use when sending a response.
// It receives the session and the triggering message so that mimic and random
// types can use live guild context.
type IdentityResolver func(s *discordgo.Session, m *discordgo.MessageCreate) (bot.Identity, error)

// BuildIdentityResolver returns the appropriate IdentityResolver for the given
// IdentityConfig. Returns an error if the config is invalid.
func BuildIdentityResolver(cfg IdentityConfig) (IdentityResolver, error) {
	switch cfg.Type {
	case "static":
		return buildStaticResolver(cfg)
	case "mimic":
		return buildMimicResolver(cfg)
	case "random":
		return buildRandomResolver(), nil
	default:
		return nil, fmt.Errorf("unknown identity type %q: must be static, mimic, or random", cfg.Type)
	}
}

func buildStaticResolver(cfg IdentityConfig) (IdentityResolver, error) {
	if cfg.BotName == "" {
		return nil, fmt.Errorf("static identity requires botName")
	}
	if cfg.AvatarURL == "" {
		return nil, fmt.Errorf("static identity requires avatarUrl")
	}
	id := bot.Identity{
		Username:  cfg.BotName,
		AvatarURL: cfg.AvatarURL,
	}
	return func(_ *discordgo.Session, _ *discordgo.MessageCreate) (bot.Identity, error) {
		return id, nil
	}, nil
}

func buildMimicResolver(cfg IdentityConfig) (IdentityResolver, error) {
	if cfg.AsMember == "" {
		return nil, fmt.Errorf("mimic identity requires as_member (Discord user ID)")
	}
	userID := cfg.AsMember
	return func(s *discordgo.Session, m *discordgo.MessageCreate) (bot.Identity, error) {
		provider := bot.NewDiscordIdentityProvider(s)
		return provider.GetIdentity(userID, m.GuildID)
	}, nil
}

func buildRandomResolver() IdentityResolver {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) (bot.Identity, error) {
		members, err := s.GuildMembers(m.GuildID, "", 1000)
		if err != nil {
			return bot.Identity{}, fmt.Errorf("random identity: could not fetch guild members: %w", err)
		}
		// Filter out bot accounts.
		nonBots := make([]*discordgo.Member, 0, len(members))
		for _, member := range members {
			if member.User != nil && !member.User.Bot {
				nonBots = append(nonBots, member)
			}
		}
		if len(nonBots) == 0 {
			return bot.Identity{}, fmt.Errorf("random identity: no non-bot guild members found")
		}
		picked := nonBots[rand.IntN(len(nonBots))]
		provider := bot.NewDiscordIdentityProvider(s)
		return provider.GetIdentity(picked.User.ID, m.GuildID)
	}
}
