package middleware

import (
	"math/rand/v2"

	"github.com/bwmarrin/discordgo"
)

// Chance passes with the given probability in [0.0, 1.0].
// A probability of 1.0 always passes; 0.0 never passes.
// Intended for low-frequency triggers (e.g. a 1% random reply).
func Chance(probability float64) MessageAuditor {
	return chanceAuditor{probability, rand.Float64}
}

type chanceAuditor struct {
	probability float64
	roll        func() float64 // injectable for testing; defaults to rand.Float64
}

func (a chanceAuditor) Audit(_ *discordgo.Session, _ *discordgo.MessageCreate) bool {
	return a.roll() < a.probability
}

