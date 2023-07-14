package command

import (
	"fmt"
	"starbunk-bot/internal/log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type HowLongTilRaid struct {
	Command    string
	CovaID     string
	RaidTeamID string
}

func (c HowLongTilRaid) CommandWord() string {
	return c.Command
}

func (c HowLongTilRaid) IsValidCommand(message string) bool {
	return isValidCommand(c.Command, message)
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
		raidTime := getNextRaid(now)

		roles, roleError := session.GuildRoles(message.GuildID)
		if roleError == nil {
			for _, role := range roles {
				log.INFO.Printf("RoleName: %s : %s", role.Name, role.ID)
			}
		}

		tag := ""
		if message.Author.ID == c.CovaID {
			tag = fmt.Sprintf("<@&%s>", c.RaidTeamID)
		}
		timeMessage := fmt.Sprintf("%s\nThe Next Raid Time is: <t:%d:f>\nWhich is <t:%d:R>", tag, raidTime.Unix(), raidTime.Unix())
		_, err := session.ChannelMessageSend(
			message.ChannelID,
			timeMessage,
		)
		if err != nil {
			log.ERROR.Println("Error calculating next raid time.")
		} else {
			session.ChannelMessageDelete(message.ChannelID, message.ChannelID)
		}
	}
}

func getNextRaid(now time.Time) time.Time {
	if now.Weekday() == time.Monday || now.Weekday() == time.Thursday {
		raidTime := time.Date(now.Year(), now.Month(), now.Day(), 24, 0, 0, 0, now.UTC().Location())
		return raidTime
	}
	return getNextRaid(now.AddDate(0, 0, 1))
}

func isTimeDST(t time.Time) bool {
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
