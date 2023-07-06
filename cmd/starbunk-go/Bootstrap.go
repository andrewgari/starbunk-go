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
	// observer.MessageService.AddObserver(reply.DeafBot{Name: "DeafBot", ID: config.UserIDs["Deaf"]})
	observer.MessageService.AddObserver(reply.EzioBot{Name: "Ezio Auditore Da Firenze", ID: config.UserIDs["Bender"]})
	observer.MessageService.AddObserver(reply.GundamBot{Name: "That Famous Unicorn Robot, \"Gandum\""})
	observer.MessageService.AddObserver(reply.HoldBot{Name: "HoldBot"})
	observer.MessageService.AddObserver(reply.MacaroniBot{Name: "MacaroniBot", ID: config.UserIDs["Venn"], Role: config.RoleIDs["Macaroni"]})
	observer.MessageService.AddObserver(reply.SheeshBot{Name: "SheeshBot"})
	observer.MessageService.AddObserver(reply.SoggyBot{Name: "SoggyBot", Role: config.RoleIDs["WetBread"]})
	observer.MessageService.AddObserver(reply.SpiderBot{Name: "Spider-Bot"})
	quoteResponses := make(map[string][]string)
	vennResponses := make([]string, 0)
	vennResponses = append(vennResponses,
		"Sorry, but that was über cringe...",
		"Geez, that was hella cringe...",
		"That was cringe to the max...",
		"What a cringe thing to say...",
		"Mondo cringe, man...",
		"Yo that was the cringiest thing I've ever heard...",
		"Your daily serving of cringe, milord...",
		"On a scale of one to cringe, that was pretty cringe...",
		"That was pretty cringe :airplane:",
		"Wow, like....cringe much?",
		"Excuse me, I seem to have dropped my cringe. Do you have it perchance?",
		"Like I always say, that was pretty cringe...",
		"C.R.I.N.G.E",
	)
	bananaSponses := make([]string, 0)
	bananaSponses = append(bananaSponses,
		"Always bring a :banana: to a party, banana's are good!",
		"Don't drop the :banana:, they're a good source of potassium!",
		"If you gave a monkey control over it's environment, it would fill the world with :banana:s...",
		"Banana. :banana:",
		"Don't judge a :banana: by it's skin.",
		"Life is full of :banana: skins.",
		"OOOOOOOOOOOOOOOOOOOOOH BA NA NA :banana:",
		":banana: Slamma!",
		"A :banana: per day keeps the Macaroni away...",
		"const bestFruit = ('b' + 'a' + + 'a').toLowerCase(); :banana:",
		"Did you know that the :banana:s we have today aren't even the same species of :banana:s we had 50 years ago. The fruit has gone extinct over time and it's actually a giant eugenics experimet to produce new species of :banana:...",
		"Monkeys always ask 'Wher :banana:', but none of them ask 'How :banana:?'",
		":banana: https://www.tiktok.com/@tracey_dintino_charles/video/7197753358143278378?_r=1&_t=8bFpt5cfIbG",
	)
	guyResponses := make([]string, 0)
	guyResponses = append(guyResponses,
		"What!? What did you say?",
		"Geeeeeet ready for Shriek Week!",
		"Try and keep up mate....",
		"But Who really died that day.\n...and who came back?",
		"Sheeeeeeeeeeeesh",
		"Rats! Rats! Weeeeeeee're the Rats!",
		"The One Piece is REEEEEEEEEEEEEEEEEEAL",
		"Psh, I dunno about that, Chief...",
		"Come to me my noble EINHERJAHR",
		"If you can't beat em, EAT em!",
		"Have you ever been so sick you sluiced your pants?",
		"Welcome back to ... Melon be Smellin'",
		"Chaotic Evil: Don't Respond. :unamused:",
		":NODDERS: Big Boys... :NODDERS:",
		"Fun Fact: That was actually in XI as well.",
		"Bird Up!",
		"Schlorp",
		"Blimbo",
	)
	log.INFO.Printf("Adding")
	quoteResponses[config.UserIDs["Guy"]] = guyResponses
	quoteResponses[config.UserIDs["Venn"]] = vennResponses
	observer.MessageService.AddObserver(reply.PickleBot{Name: "GremlinBot", ID: config.UserIDs["Sig"]})
	observer.MessageService.AddObserver(reply.VennBot{GuildID: config.GuildIDs["Starbunk"], UserID: config.UserIDs["Venn"], Responses: vennResponses, Bananasponses: bananaSponses})
	observer.MessageService.AddObserver(reply.GuyBot{GuildID: config.GuildIDs["Starbunk"], UserID: config.UserIDs["Guy"], Responses: guyResponses})
	observer.MessageService.AddObserver(command.MusicCorrect{})
	observer.MessageService.AddObserver(
		reply.RagtimeBot{
			TriviaMaster:        config.RoleIDs["TriviaMaster"],
			TriviaChannel:       config.ChannelIDs["Trivia"],
			TriviaReviewChannel: config.ChannelIDs["TriviaReview"],
		})
	observer.MessageService.AddObserver(reply.SixtyNineBot{Name: "CovaBot"})
}

func RegisterCommandBots() {
	observer.CommandBots["clearWebhooks"] = command.ClearWebhooks{Command: "clearWebhooks", GuildID: config.GuildIDs["Starbunk"]}
	observer.CommandBots["raidwhen"] =
		command.HowLongTilRaid{Command: "raidwhen"}
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
