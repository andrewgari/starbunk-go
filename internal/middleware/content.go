package middleware

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// HasContent drops messages with empty or whitespace-only content.
var HasContent MessageAuditor = hasContentAuditor{}

// ContentContains passes only messages whose content includes substr.
func ContentContains(substr string) MessageAuditor {
	return contentContainsAuditor{substr}
}

// ContentMatches passes only messages whose content matches the given pattern.
func ContentMatches(re *regexp.Regexp) MessageAuditor {
	return contentMatchesAuditor{re}
}

// HasAttachment passes only messages that include at least one file attachment.
var HasAttachment MessageAuditor = hasAttachmentAuditor{}

// — implementations —

type hasContentAuditor struct{}

func (hasContentAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
	return strings.TrimSpace(m.Content) != ""
}

type contentContainsAuditor struct{ substr string }

func (a contentContainsAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
	return strings.Contains(m.Content, a.substr)
}

type contentMatchesAuditor struct{ re *regexp.Regexp }

func (a contentMatchesAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
	return a.re.MatchString(m.Content)
}

type hasAttachmentAuditor struct{}

func (hasAttachmentAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
	return len(m.Attachments) > 0
}
