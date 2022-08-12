package main

import (
	"starbunk-bot/internal/bot/command"
	"starbunk-bot/internal/bot/reply"
	"starbunk-bot/internal/bot/voice"
	"starbunk-bot/internal/config"
	"starbunk-bot/internal/log"
	"starbunk-bot/internal/observer"
)

func RegisterReplyBots() {
	observer.MessageService.AddObserver(reply.BluBot{Name: "BluBot"})
	observer.MessageService.AddObserver(reply.ChaosBot{Name: "ChaosBot"})
	observer.MessageService.AddObserver(reply.CheckBot{Name: "CzechBot"})
	observer.MessageService.AddObserver(reply.DeafBot{Name: "DeafBot", ID: config.UserIDs["Deaf"]})
	observer.MessageService.AddObserver(reply.EzioBot{Name: "Ezio Auditore Da Firenze", ID: config.UserIDs["Bender"]})
	observer.MessageService.AddObserver(reply.GundamBot{Name: "That Famous Unicorn Robot, \"Gandum\""})
	observer.MessageService.AddObserver(reply.HoldBot{Name: "HoldBot"})
	observer.MessageService.AddObserver(reply.MacaroniBot{Name: "MacaroniBot", ID: config.UserIDs["Venn"], Role: config.RoleIDs["Macaroni"]})
	observer.MessageService.AddObserver(reply.PickleBot{Name: "GremlinBot", ID: config.UserIDs["Sig"]})
	observer.MessageService.AddObserver(reply.SheeshBot{Name: "SheeshBot"})
	observer.MessageService.AddObserver(reply.SixtyNineBot{Name: "CovaBot"})
	observer.MessageService.AddObserver(reply.SoggyBot{Name: "SoggyBot", Role: config.RoleIDs["WetBread"]})
	observer.MessageService.AddObserver(reply.SpiderBot{Name: "Spider-Bot"})
	observer.MessageService.AddObserver(reply.VennBot{ID: config.UserIDs["Venn"]})
}

func RegisterCommandBots() {
	observer.CommandBots["clearWebhooks"] = command.ClearWebhooks{Command: "clearWebhooks", GuildID: config.GuildIDs["Starbunk"]}
	observer.CommandBots["play"] = command.MusicCorrect{Command: "play"}
	observer.CommandBots["covaPlay"] = command.MusicCorrect{Command: "covaPlay"}
}

func RegisterVoiceBots() {
	log.WARN.Println("Adding Voice Bots")
	var guyBot observer.IVoiceObserver = voice.GuyChannelBot{
		GuysID:           config.UserIDs["Guy"],
		GuysChannelID:    config.ChannelIDs["OnlyGuy"],
		NotGuysChannelID: config.ChannelIDs["NoGuy"],
		LoungeId:         config.ChannelIDs["Lounge"],
		GuildID:          config.GuildIDs["Starbunk"],
	}
	type GuyChannelBot struct {
		GuysID           string
		GuysChannelID    string
		NotGuysChannelID string
		LoungeId         string
		GuildID          string
	}
	observer.VoiceService.AddObserver(guyBot)
	var feliBot observer.IVoiceObserver = voice.FeliBot{
		Name:            "FeliBot",
		FeliID:          config.UserIDs["Feli"],
		GuildID:         config.GuildIDs["Starbunk"],
		AFK_ID:          config.ChannelIDs["AFK"],
		WhaleWatchersID: config.ChannelIDs["WhaleWatchers"],
	}
	observer.VoiceService.AddObserver(feliBot)
}
