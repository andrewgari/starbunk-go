package utils

import (
	"golang-discord-bot/internal/log"
	"regexp"
)

func Match(pattern string, s string) bool {
	match, err := regexp.MatchString(pattern, s)
	if err != nil {
		log.ERROR.Println("Error Parsing Message: ", err)
	}
	return match
}
