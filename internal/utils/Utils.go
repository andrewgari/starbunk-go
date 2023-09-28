package utils

import (
	"math/rand"
	"regexp"
	"starbunk-bot/internal/log"
	"time"

	"github.com/bwmarrin/discordgo"
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

func getGuildMember(session *discordgo.Session, guildID string, userID string) (discordgo.Member, error) {
	member, err := session.GuildMember(guildID, userID)
	return *member, err
}

func GetNickname(session *discordgo.Session, message discordgo.Message) string {
	member, err := getGuildMember(session, message.GuildID, message.Author.ID)
	if err != nil {
		return message.Author.Username
	}
	return member.Nick
}

func GetAvatarUrl(session *discordgo.Session, message discordgo.Message) string {
	member, err := getGuildMember(session, message.GuildID, message.Author.ID)
	if err != nil {
		return message.Author.AvatarURL("")
	}
	return member.AvatarURL("")
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
