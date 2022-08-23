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
	rand.Seed(time.Now().UnixNano())
	var roll = rand.Intn(100)
	log.INFO.Printf("Rolled a Percent Chance: %d", roll)
	return roll < target
}
