package replybot

import (
	"math/rand"
	"time"
)

func roll20(target int) bool {
	rand.Seed(time.Now().UnixNano())
	var roll = rand.Intn(20-1) + 1
	return roll >= target
}
