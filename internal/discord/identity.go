package discord

import (
	"github.com/bwmarrin/discordgo"
)

// Identity represents the persona a bot or poster assumes when sending messages.
type Identity struct {
	Nickname  string
	Username  string
	AvatarURL string
	Metadata  map[string]string
}

// IsValid reports whether the Identity carries the minimum fields required to
// execute a webhook: a display name and avatar URL.
func (i Identity) IsValid() bool {
	return i.Username != "" && i.AvatarURL != ""
}

// Resolve returns a complete Identity, falling back to the bot's own profile
// from the Discord session for any missing Username or AvatarURL.
func (i Identity) Resolve(s *discordgo.Session) Identity {
	resolved := i

	if s != nil && s.State != nil && s.State.User != nil {
		if resolved.Username == "" {
			resolved.Username = s.State.User.Username
		}
		if resolved.AvatarURL == "" {
			resolved.AvatarURL = s.State.User.AvatarURL("")
		}
	}

	if resolved.Metadata == nil {
		resolved.Metadata = make(map[string]string)
	}

	return resolved
}

// IdentityProvider retrieves a Discord user's identity on demand.
type IdentityProvider interface {
	GetIdentity(userID string, guildID string) (Identity, error)
}

// DiscordIdentityProvider resolves user identities directly from Discord,
// preferring guild-member details (nick, server avatar) over global user details.
type DiscordIdentityProvider struct {
	session *discordgo.Session
}

// NewDiscordIdentityProvider creates a DiscordIdentityProvider backed by s.
func NewDiscordIdentityProvider(s *discordgo.Session) *DiscordIdentityProvider {
	return &DiscordIdentityProvider{session: s}
}

// GetIdentity queries Discord for userID. When guildID is non-empty it prefers
// the guild member record (nick, server avatar) over the global user profile.
func (p *DiscordIdentityProvider) GetIdentity(userID string, guildID string) (Identity, error) {
	var id Identity

	if guildID != "" {
		member, err := p.session.GuildMember(guildID, userID)
		if err == nil && member != nil {
			id.Username = member.User.Username
			id.Nickname = member.Nick
			id.AvatarURL = member.AvatarURL("")
			if id.AvatarURL == "" && member.User != nil {
				id.AvatarURL = member.User.AvatarURL("")
			}
			return id, nil
		}
	}

	user, err := p.session.User(userID)
	if err != nil {
		return id, err
	}

	id.Username = user.Username
	id.AvatarURL = user.AvatarURL("")

	return id, nil
}
