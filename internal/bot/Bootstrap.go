package bot

import (
	"starbunk-bot/internal/bot/command"
	"starbunk-bot/internal/bot/reply"
	"starbunk-bot/internal/config"
	"starbunk-bot/internal/observer"
)

var CommandBots = make(map[string]command.ICommandBot)

func RegisterReplyBots() {
	var bluBot observer.IMessageObserver = reply.BluBot{Name: "BluBot"}
	observer.MessageService.AddObserver(bluBot)
	var chaosBot observer.IMessageObserver = reply.ChaosBot{Name: "ChaosBot"}
	observer.MessageService.AddObserver(chaosBot)
	var checkBot observer.IMessageObserver = reply.CheckBot{Name: "CzechBot"}
	observer.MessageService.AddObserver(checkBot)
	var deafBot observer.IMessageObserver = reply.DeafBot{Name: "DeafBot", ID: config.UserIDs["Deaf"]}
	observer.MessageService.AddObserver(deafBot)
	var ezioBot observer.IMessageObserver = reply.EzioBot{Name: "Ezio Auditore Da Firenze", ID: config.UserIDs["Bender"]}
	observer.MessageService.AddObserver(ezioBot)
	var gundamBot observer.IMessageObserver = reply.GundamBot{Name: "That Famous Unicorn Robot, \"Gandum\""}
	observer.MessageService.AddObserver(gundamBot)
	var holdBot observer.IMessageObserver = reply.HoldBot{Name: "HoldBot"}
	observer.MessageService.AddObserver(holdBot)
	var macaroniBot observer.IMessageObserver = reply.MacaroniBot{Name: "MacaroniBot", ID: config.UserIDs["Venn"], Role: config.RoleIDs["Macaroni"]}
	observer.MessageService.AddObserver(macaroniBot)
	var pickleBot observer.IMessageObserver = reply.PickleBot{Name: "GremlinBot", ID: config.UserIDs["Sig"]}
	observer.MessageService.AddObserver(pickleBot)
	var sheeshBot observer.IMessageObserver = reply.SheeshBot{Name: "SheeshBot", ID: config.UserIDs["Guy"]}
	observer.MessageService.AddObserver(sheeshBot)
	var sixtyNineBot observer.IMessageObserver = reply.SixtyNineBot{Name: "CovaBot"}
	observer.MessageService.AddObserver(sixtyNineBot)
	var soggyBot observer.IMessageObserver = reply.SoggyBot{Name: "SoggyBot", Role: config.RoleIDs["WetBread"]}
	observer.MessageService.AddObserver(soggyBot)
	var spiderBot observer.IMessageObserver = reply.SpiderBot{Name: "Spider-Bot"}
	observer.MessageService.AddObserver(spiderBot)
	var vennBot observer.IMessageObserver = reply.VennBot{ID: config.UserIDs["Venn"]}
	observer.MessageService.AddObserver(vennBot)
}

func RegisterCommandBots() {
	var clearWebhooks command.ICommandBot = command.ClearWebhooks{Command: "clearWebhooks", GuildID: config.GuildIDs["BLU"]}
	CommandBots["clearWebhooks"] = clearWebhooks
}
