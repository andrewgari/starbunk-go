package replybot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andrewgari/starbunk-go/internal/middleware"
	"github.com/bwmarrin/discordgo"
)

// Compile converts a ConditionNode tree into a single MessageAuditor.
// Combinators are checked first (AllOf, AnyOf, NoneOf), then leaf conditions.
// Returns an error if the node is empty or if a regex pattern is invalid.
func Compile(node ConditionNode) (middleware.MessageAuditor, error) {
	// Combinators take precedence over leaves.
	if len(node.AllOf) > 0 {
		children, err := compileAll(node.AllOf)
		if err != nil {
			return nil, fmt.Errorf("all_of: %w", err)
		}
		return middleware.AllOf(children...), nil
	}
	if len(node.AnyOf) > 0 {
		children, err := compileAll(node.AnyOf)
		if err != nil {
			return nil, fmt.Errorf("any_of: %w", err)
		}
		return middleware.AnyOf(children...), nil
	}
	if len(node.NoneOf) > 0 {
		children, err := compileAll(node.NoneOf)
		if err != nil {
			return nil, fmt.Errorf("none_of: %w", err)
		}
		return middleware.Not(middleware.AnyOf(children...)), nil
	}

	// Leaf conditions — checked in a stable priority order.
	if node.ContainsWord != "" {
		pattern := fmt.Sprintf(`(?i)\b%s\b`, regexp.QuoteMeta(node.ContainsWord))
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("contains_word %q: %w", node.ContainsWord, err)
		}
		return containsWordAuditor{re}, nil
	}
	if node.ContainsPhrase != "" {
		return containsPhraseAuditor{strings.ToLower(node.ContainsPhrase)}, nil
	}
	if node.MatchesPattern != "" {
		re, err := regexp.Compile(node.MatchesPattern)
		if err != nil {
			return nil, fmt.Errorf("matches_pattern %q: %w", node.MatchesPattern, err)
		}
		return middleware.ContentMatches(re), nil
	}
	if node.FromUser != "" {
		return middleware.AuthorID(node.FromUser), nil
	}
	if node.WithChance != nil {
		return middleware.Chance(float64(*node.WithChance) / 100), nil
	}
	if node.Always {
		return alwaysAuditor{}, nil
	}

	return nil, fmt.Errorf("empty condition node: at least one condition must be set")
}

// compileAll compiles a slice of ConditionNodes and returns the auditors.
func compileAll(nodes []ConditionNode) ([]middleware.MessageAuditor, error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("combinator requires at least one child condition")
	}
	auditors := make([]middleware.MessageAuditor, 0, len(nodes))
	for i, n := range nodes {
		a, err := Compile(n)
		if err != nil {
			return nil, fmt.Errorf("child[%d]: %w", i, err)
		}
		auditors = append(auditors, a)
	}
	return auditors, nil
}

// containsWordAuditor matches whole words case-insensitively using \bword\b.
type containsWordAuditor struct{ re *regexp.Regexp }

func (a containsWordAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
	return a.re.MatchString(m.Content)
}

// containsPhraseAuditor matches a phrase as a case-insensitive substring.
type containsPhraseAuditor struct{ phrase string }

func (a containsPhraseAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
	return strings.Contains(strings.ToLower(m.Content), a.phrase)
}

// alwaysAuditor passes every message unconditionally.
type alwaysAuditor struct{}

func (alwaysAuditor) Audit(_ *discordgo.Session, _ *discordgo.MessageCreate) bool {
	return true
}

// neverAuditor rejects every message (used for invalid ignore_bots+ignore_humans combo).
type neverAuditor struct{}

func (neverAuditor) Audit(_ *discordgo.Session, _ *discordgo.MessageCreate) bool {
	return false
}
