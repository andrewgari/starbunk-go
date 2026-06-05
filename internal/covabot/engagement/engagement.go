package engagement

import "sync"

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
	Muted         bool
	Dampened      bool
	LastCovaSpeak bool
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
	lastSpeak := state.LastCovaSpeak
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

	// 5. Engagement Continuity
	if lastSpeak {
		return Result{
			Respond: true,
			Reason:  ReasonContext,
			Energy:  EnergyNormal,
		}
	}

	return Result{Respond: false}
}

// RecordCovaSpeak records that CovaBot just spoke in a channel, increasing engagement continuity.
func (m *Manager) RecordCovaSpeak(channelID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.getState(channelID).LastCovaSpeak = true
}

// Dampen temporarily raises the pull floor in a channel, silencing non-directed responses.
func (m *Manager) Dampen(channelID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.getState(channelID).Dampened = true
}

// SetMute applies a hard floor. Only direct addresses pass through.
func (m *Manager) SetMute(channelID string, mute bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.getState(channelID).Muted = mute
}
