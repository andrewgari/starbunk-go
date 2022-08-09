package utils

import (
	"golang-discord-bot/internal/log"
	"math/rand"
	"regexp"
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

func Roll20(target int) bool {
	rand.Seed(time.Now().UnixNano())
	var roll = rand.Intn(20-1) + 1
	return roll >= target
}
