package middleware

import (
	"slices"

	"github.com/bwmarrin/discordgo"
)

// NotSelf drops messages sent by the bot itself.
var NotSelf MessageAuditor = notSelfAuditor{}

// NotBot drops messages where the author is any bot account.
var NotBot MessageAuditor = notBotAuditor{}

// IsBot passes only messages where the author is a bot account.
var IsBot MessageAuditor = Not(NotBot)

// AuthorID passes only messages from the given Discord user ID.
func AuthorID(id string) MessageAuditor {
	return authorIDAuditor{id}
}

// NotAuthorID drops messages from the given Discord user ID.
func NotAuthorID(id string) MessageAuditor {
	return Not(AuthorID(id))
}

// AuthorNamed passes only messages whose author's username equals name
// (case-sensitive, matches Username not display name).
func AuthorNamed(name string) MessageAuditor {
	return authorNamedAuditor{name}
}

// AuthorHasRole passes only messages where the author holds a given role ID in
// the guild. Drops the message if guild member data is unavailable.
func AuthorHasRole(roleID string) MessageAuditor {
	return authorHasRoleAuditor{roleID}
}

// — implementations —

type notSelfAuditor struct{}

func (notSelfAuditor) Audit(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	return s.State.User == nil || m.Author.ID != s.State.User.ID
}

type notBotAuditor struct{}

func (notBotAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
	return !m.Author.Bot
}

type authorIDAuditor struct{ id string }

func (a authorIDAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
	return m.Author.ID == a.id
}

type authorNamedAuditor struct{ name string }

func (a authorNamedAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
	return m.Author.Username == a.name
}

type authorHasRoleAuditor struct{ roleID string }

func (a authorHasRoleAuditor) Audit(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	member, err := s.State.Member(m.GuildID, m.Author.ID)
	if err != nil {
		return false
	}
	return slices.Contains(member.Roles, a.roleID)
}
