package command

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"starbunk-bot/internal/log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type CreateEventBot struct {
	Command             string
	NebulaLeadRole      string
	NebulaBunker        string
	NebulaAnnouncements string
}

var eventName string = "Nebula Raid"

func (c CreateEventBot) CommandWord() string {
	return c.Command
}

func (c CreateEventBot) IsValidCommand(message string) bool {
	return isValidCommand(c.Command, message)
}

func (c CreateEventBot) hasPermissions(message discordgo.Message) bool {
	return c.contains(message.Member.Roles, c.NebulaLeadRole)
}

func (c CreateEventBot) contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func (c CreateEventBot) ProcessMessage(session *discordgo.Session, message discordgo.Message) {
	if c.IsValidCommand(message.Content) && c.hasPermissions(message) {
		log.INFO.Println("Valid Command")
		c.createRaidEvent(session, message.GuildID, time.Thursday)
		c.createRaidEvent(session, message.GuildID, time.Sunday)
		events, err := session.GuildScheduledEvents(message.GuildID, true)
		if err != nil {
			log.ERROR.Println("Error Getting Scheduled Events")
			return
		}
		for _, event := range events {
			// index is the index where we are
			// element is the element from someSlice for where we are
			if event.Name == eventName {
				// https://discord.com/events/753251582719688714/1017153913503363234
				session.ChannelMessageSend(c.NebulaAnnouncements, fmt.Sprintf("https://discord.com/events/%v/%v", event.GuildID, event.ID))
			}
		}
	} else {
		log.WARN.Println("Invalid Command: ", message.Content, c.hasPermissions(message), c.IsValidCommand(message.Content))
	}
}

func (c CreateEventBot) createRaidEvent(session *discordgo.Session, guildId string, weekday time.Weekday) {
	now := time.Now()
	startHour := 0
	startMinute := 0

	endHour := 3
	endMinute := 0

	switch weekday {
	case time.Thursday:
		startHour = 1
		startMinute = 30
	case time.Sunday:
		startHour = 0
		startMinute = 0
	default:
		log.WARN.Println("Error not set to Wednesday or Saturday")
		return
	}

	getFileHash()

	raidStart := nextEvent(now, weekday, startHour, startMinute)
	raidEnd := nextEvent(now, weekday, endHour, endMinute)
	scheduledEvent, err := session.GuildScheduledEventCreate(guildId, &discordgo.GuildScheduledEventParams{
		Name:               eventName,
		Description:        "The important thing is we tried...",
		ScheduledStartTime: &raidStart,
		ScheduledEndTime:   &raidEnd,
		EntityType:         discordgo.GuildScheduledEventEntityTypeVoice,
		ChannelID:          c.NebulaBunker,
		PrivacyLevel:       discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
		Image:              getFileHash(),
	})
	if err != nil {
		log.ERROR.Printf("Error creating scheduled event: %v", err)
		return
	}
	log.INFO.Println("Created Scheduled Event: ", scheduledEvent.Name)
}

func nextEvent(t time.Time, weekday time.Weekday, hour, minute int) time.Time {
	days := int((7 + (weekday - t.Weekday())) % 7)
	year, month, day := t.AddDate(0, 0, days).Date()
	return time.Date(year, month, day, hour, minute, 0, 0, time.UTC)
}

func getFileHash() string {
	bytes, err := ioutil.ReadFile("resources/nebula.webp")
	if err != nil {
		log.ERROR.Println("Error Reading File: ", err)
	}

	var base64Encoding string
	// Determine the content type of the image file
	mimeType := http.DetectContentType(bytes)

	// Prepend the appropriate URI scheme header depending
	// on the MIME type
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	case "image/webp":
		base64Encoding += "data:image/webp;base64,"
	}

	// Append the base64 encoded output
	base64Encoding += base64.StdEncoding.EncodeToString(bytes)
	return base64Encoding
}
