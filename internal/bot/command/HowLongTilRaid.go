package command

import (
	"fmt"
	"starbunk-bot/internal/log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type HowLongTilRaid struct {
	Command string
}

func (c HowLongTilRaid) CommandWord() string {
	return c.Command
}

func (c HowLongTilRaid) IsValidCommand(message string) bool {
	return isValidCommand(c.Command, message)
}

func (c HowLongTilRaid) WednesdayRaidTime() {

}

func (c HowLongTilRaid) SaturdayRaidTime() {
	now := time.Now()
	now.Month()

}

func (c HowLongTilRaid) ProcessMessage(session *discordgo.Session, message discordgo.Message) {
	if c.IsValidCommand(message.Content) {

		pstLoc, timeError := time.LoadLocation("America/Los_Angeles")
		if timeError != nil {
			fmt.Println(timeError)
			return
		}
		utc := time.Now().UTC()
		now := utc.In(pstLoc)
		fmt.Println(now, now.Location(), ": DST", isTimeDST(now))

		raidTime := getNextRaid(now)
		log.INFO.Println(raidTime)
		diff := raidTime.Sub(now)
		log.INFO.Println(diff)
		seconds := diff.Seconds()
		log.INFO.Println(seconds)
		days := 0
		for seconds >= 86400 {
			days++
			seconds -= 86400
		}
		log.INFO.Println(days)
		hours := 0
		for seconds >= 3600 {
			hours++
			seconds -= 3600
		}
		log.INFO.Println(hours)
		minutes := 0
		for seconds >= 60 {
			minutes++
			seconds -= 60
		}
		log.INFO.Println(minutes)
		if now.IsDST() {
			hours -= 1
		}
		_, err := session.ChannelMessageSend(
			message.ChannelID,
			fmt.Sprintf("Raid is in %d days, %d hours and %d minutes", days, hours, minutes),
		)
		if err != nil {
			log.ERROR.Println("Error calculating next raid time.")
		}
	}
}

func getNextRaid(now time.Time) time.Time {
	switch now.Weekday() {
	case time.Thursday:
		raidTime := time.Date(now.Year(), now.Month(), now.Day(), 2, 30, 0, 0, now.UTC().Location())
		return raidTime
	case time.Sunday:
		raidTime := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, now.UTC().Location())
		return raidTime
	default:
		return getNextRaid(now.AddDate(0, 0, 1))
	}
}

func isTimeDST(t time.Time) bool {
	// If the most recent (within the last year) clock change
	// was forward then assume the change was from std to dst.
	hh, mm, _ := t.UTC().Clock()
	tClock := hh*60 + mm
	for m := -1; m > -12; m-- {
		// assume dst lasts for least one month
		hh, mm, _ := t.AddDate(0, m, 0).UTC().Clock()
		clock := hh*60 + mm
		if clock != tClock {
			log.INFO.Println("It is DST")
			return clock > tClock
		}
	}
	log.INFO.Println("It is not DST")
	// assume no dst
	return false
}
