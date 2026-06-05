package engagement

import (
	"sync"
	"time"
)

// recentSpeakWindow is how long engagement continuity stays active after Cova speaks.
// After this window the flag decays and Cova stops following ambient channel messages.
const recentSpeakWindow = 5 * time.Minute

// MessageInput represents the minimal information needed from a message
// to judge engagement pull and apply restraint.
type MessageInput struct {
	ChannelID       string
	AuthorID        string
	IsMentioned     bool
	IsReplyToMe     bool
	IsAddresseeSelf bool
}

// GateReason articulates the intent of the response.
type GateReason string

const (
	ReasonMention GateReason = "direct_mention"
	ReasonReply   GateReason = "reply_to_cova"
	ReasonContext GateReason = "engagement_continuity"
)

// GateEnergy drives the cadence and length of the response.
type GateEnergy string

const (
	EnergyQuickJab GateEnergy = "quick-jab"
	EnergyNormal   GateEnergy = "normal"
	EnergyInvested GateEnergy = "invested"
)

// Result is the outcome of the ShouldRespond check.
type Result struct {
	Respond bool
	Reason  GateReason
	Energy  GateEnergy
}

type channelState struct {
	Muted       bool
	Dampened    bool
	lastSpokeAt time.Time // zero value = Cova has not spoken; decays after recentSpeakWindow
}

// Manager tracks the engagement state per channel and decides if CovaBot should respond.
type Manager struct {
	mu     sync.Mutex
	states map[string]*channelState
}

// NewManager creates a new engagement manager.
func NewManager() *Manager {
	return &Manager{
		states: make(map[string]*channelState),
	}
}

func (m *Manager) getState(channelID string) *channelState {
	state, ok := m.states[channelID]
	if !ok {
		state = &channelState{}
		m.states[channelID] = state
	}
	return state
}

// ShouldRespond decides if CovaBot should respond to the incoming message.
func (m *Manager) ShouldRespond(input MessageInput) Result {
	m.mu.Lock()
	state := m.getState(input.ChannelID)
	muted := state.Muted
	dampened := state.Dampened
	recentlySpoke := !state.lastSpokeAt.IsZero() && time.Since(state.lastSpokeAt) < recentSpeakWindow
	m.mu.Unlock()

	// 1. Direct Mention (Highest pull, clears all restraints including mute)
	if input.IsMentioned {
		return Result{
			Respond: true,
			Reason:  ReasonMention,
			Energy:  EnergyInvested,
		}
	}

	// 2. Mute stops everything except direct mention
	if muted {
		return Result{Respond: false}
	}

	// 3. Direct Reply or Addressee==Self (High pull, clears dampener)
	if input.IsReplyToMe || input.IsAddresseeSelf {
		return Result{
			Respond: true,
			Reason:  ReasonReply,
			Energy:  EnergyNormal,
		}
	}

	// 4. Dampener stops ambient/continuity
	if dampened {
		return Result{Respond: false}
	}

	// 5. Engagement Continuity — active for recentSpeakWindow after Cova last spoke
	if recentlySpoke {
		return Result{
			Respond: true,
			Reason:  ReasonContext,
			Energy:  EnergyNormal,
		}
	}

	return Result{Respond: false}
}

// RecordCovaSpeak records that CovaBot just spoke in a channel.
// Engagement continuity will remain active for recentSpeakWindow from this call.
func (m *Manager) RecordCovaSpeak(channelID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.getState(channelID).lastSpokeAt = time.Now()
}

// Dampen temporarily raises the pull floor in a channel, silencing non-directed responses.
// NOTE: auto-decay (design doc §4.3) is not yet implemented — the dampener currently
// stays active until explicitly cleared. Use SetDampen(channelID, false) to lift it.
func (m *Manager) Dampen(channelID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.getState(channelID).Dampened = true
}

// SetDampen explicitly sets or clears the dampener for a channel.
func (m *Manager) SetDampen(channelID string, dampened bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.getState(channelID).Dampened = dampened
}

// SetMute applies a hard floor. Only direct addresses pass through.
func (m *Manager) SetMute(channelID string, mute bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.getState(channelID).Muted = mute
}
