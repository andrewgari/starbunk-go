// Package replybot implements the YAML-driven reply-bot system for BunkBot.
// A reply bot is defined in a YAML config file and consists of a bot identity,
// a set of message triggers (each with a condition tree and response pool), and
// flags that control which message authors are eligible.
package replybot

// FileConfig is the top-level structure of a YAML config file.
type FileConfig struct {
	ReplyBots []BotConfig `yaml:"reply-bots"`
}

// BotConfig is one entry under reply-bots.
type BotConfig struct {
	Name         string          `yaml:"name"`
	Identity     IdentityConfig  `yaml:"identity"`
	Responses    []string        `yaml:"responses"`
	Triggers     []TriggerConfig `yaml:"triggers"`
	IgnoreBots   bool            `yaml:"ignore_bots"`
	IgnoreHumans bool            `yaml:"ignore_humans"`
}

// IdentityConfig describes how the bot's persona is resolved at send time.
// Type must be one of: "static", "mimic", "random".
type IdentityConfig struct {
	Type      string `yaml:"type"`
	BotName   string `yaml:"botName"`   // required for type=static
	AvatarURL string `yaml:"avatarUrl"` // required for type=static
	AsMember  string `yaml:"as_member"` // required for type=mimic (Discord user ID)
}

// TriggerConfig is one trigger entry within a bot.
type TriggerConfig struct {
	Name       string        `yaml:"name"`
	Conditions ConditionNode `yaml:"conditions"`
	Responses  []string      `yaml:"responses"`
}

// ConditionNode is a node in the recursive condition logic tree.
// Exactly one field should be set per node.
//
// Leaf conditions: ContainsWord, ContainsPhrase, MatchesPattern, FromUser, WithChance, Always.
// Combinators: AllOf (AND), AnyOf (OR), NoneOf (NOR).
//
// WithChance uses a pointer so that with_chance: 0 (never fire) is distinguishable
// from the field being absent.
type ConditionNode struct {
	// Leaf conditions
	ContainsWord   string `yaml:"contains_word"`
	ContainsPhrase string `yaml:"contains_phrase"`
	MatchesPattern string `yaml:"matches_pattern"`
	FromUser       string `yaml:"from_user"`
	WithChance     *int   `yaml:"with_chance"` // 0–100; pointer so 0 is valid
	Always         bool   `yaml:"always"`

	// Logical combinators
	AllOf  []ConditionNode `yaml:"all_of"`
	AnyOf  []ConditionNode `yaml:"any_of"`
	NoneOf []ConditionNode `yaml:"none_of"`
}

// isEmpty reports whether no condition has been set on this node.
func (n ConditionNode) isEmpty() bool {
	return n.ContainsWord == "" &&
		n.ContainsPhrase == "" &&
		n.MatchesPattern == "" &&
		n.FromUser == "" &&
		n.WithChance == nil &&
		!n.Always &&
		len(n.AllOf) == 0 &&
		len(n.AnyOf) == 0 &&
		len(n.NoneOf) == 0
}
