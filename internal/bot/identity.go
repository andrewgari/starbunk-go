package bot

import (
	"github.com/bwmarrin/discordgo"
)

// Identity represents the persona a bot or poster assumes.
type Identity struct {
	Nickname  string
	Username  string
	AvatarURL string
	Metadata  map[string]string
}

// IsValid checks if the Identity has the required username and avatar url.
func (i Identity) IsValid() bool {
	return i.Username != "" && i.AvatarURL != ""
}

// Resolve returns a complete Identity, falling back to the bot's default profile
// name and avatar from the Discord session if they are not specified.
func (i Identity) Resolve(s *discordgo.Session) Identity {
	resolved := i

	// If Username or AvatarURL is missing, use the bot's default profile
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

// IdentityProvider defines how we retrieve identities for a given user.
type IdentityProvider interface {
	GetIdentity(userID string, guildID string) (Identity, error)
}

// DiscordIdentityProvider retrieves user identities directly from Discord.
type DiscordIdentityProvider struct {
	session *discordgo.Session
}

// NewDiscordIdentityProvider creates a new DiscordIdentityProvider.
func NewDiscordIdentityProvider(s *discordgo.Session) *DiscordIdentityProvider {
	return &DiscordIdentityProvider{session: s}
}

// GetIdentity queries Discord for a user's identity, preferring server-specific details (Member) if a guildID is provided.
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
