// Package middleware provides the MessageAuditor interface and composable gate
// primitives for evaluating Discord messages before they reach a bot handler.
//
// Every gate satisfies the same MessageAuditor interface, so any combination of
// gates can be composed freely using the AllOf, AnyOf, and Not combinators:
//
//	middleware.AllOf(
//	    middleware.IsBot,
//	    middleware.AuthorNamed("Jeff"),
//	    middleware.Not(middleware.OnWeekdays(time.Friday, time.Tuesday)),
//	)
//
// New gates are added by implementing MessageAuditor and placing them in the
// appropriate file (author.go, content.go, context.go) or a new file for a new
// category.
package middleware

import "github.com/bwmarrin/discordgo"

// MessageAuditor is the mandatory evaluation gate for incoming Discord messages.
// bot.Run requires one. The framework automatically applies it to every
// MessageCreate event before invoking any registered handler.
//
// Audit returns true if the message should be processed, false if it should
// be dropped silently.
type MessageAuditor interface {
	Audit(s *discordgo.Session, m *discordgo.MessageCreate) bool
}

// AllOf returns an auditor that passes only when every child passes.
// Short-circuits on the first failure. Passes vacuously when given no children.
func AllOf(auditors ...MessageAuditor) MessageAuditor {
	return allOfAuditor{auditors}
}

// AnyOf returns an auditor that passes when at least one child passes.
// Short-circuits on the first success. Fails vacuously when given no children.
func AnyOf(auditors ...MessageAuditor) MessageAuditor {
	return anyOfAuditor{auditors}
}

// Not inverts an auditor.
func Not(a MessageAuditor) MessageAuditor {
	return notAuditor{a}
}

type allOfAuditor struct{ auditors []MessageAuditor }

func (a allOfAuditor) Audit(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	for _, child := range a.auditors {
		if !child.Audit(s, m) {
			return false
		}
	}
	return true
}

type anyOfAuditor struct{ auditors []MessageAuditor }

func (a anyOfAuditor) Audit(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	for _, child := range a.auditors {
		if child.Audit(s, m) {
			return true
		}
	}
	return false
}

type notAuditor struct{ inner MessageAuditor }

func (a notAuditor) Audit(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	return !a.inner.Audit(s, m)
}
