package middleware

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// GuildOnly drops direct messages (messages with no GuildID).
var GuildOnly MessageAuditor = guildOnlyAuditor{}

// DMOnly drops guild messages, passing only direct messages.
var DMOnly MessageAuditor = dmOnlyAuditor{}

// InChannel passes only messages sent in the given channel ID.
func InChannel(channelID string) MessageAuditor {
	return inChannelAuditor{channelID}
}

// OnWeekdays passes only messages sent on one of the given weekdays (in UTC).
// Example: OnWeekdays(time.Monday, time.Wednesday, time.Friday)
func OnWeekdays(days ...time.Weekday) MessageAuditor {
	set := make(map[time.Weekday]struct{}, len(days))
	for _, d := range days {
		set[d] = struct{}{}
	}
	return onWeekdaysAuditor{set}
}

// — implementations —

type guildOnlyAuditor struct{}

func (guildOnlyAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
	return m.GuildID != ""
}

type dmOnlyAuditor struct{}

func (dmOnlyAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
	return m.GuildID == ""
}

type inChannelAuditor struct{ channelID string }

func (a inChannelAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
	return m.ChannelID == a.channelID
}

type onWeekdaysAuditor struct{ days map[time.Weekday]struct{} }

func (a onWeekdaysAuditor) Audit(_ *discordgo.Session, m *discordgo.MessageCreate) bool {
	_, ok := a.days[m.Timestamp.UTC().Weekday()]
	return ok
}
