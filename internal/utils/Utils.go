package utils

import (
	"math/rand"
	"regexp"
	"starbunk-bot/internal/log"
	"time"
)

func Match(pattern string, s string) bool {
	match, err := regexp.MatchString("(?i)"+pattern, s)
	if err != nil {
		log.ERROR.Println("Error Parsing Message: ", err)
		return false
	}
	return match
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func PercentChance(target int) bool {
	roll := RandomRoll(100)
	log.INFO.Printf("Rolled a Percent Chance: %d", roll)
	return roll < target
}

func RandomRoll(limit int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(limit)
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
