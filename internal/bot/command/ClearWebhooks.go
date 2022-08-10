package command

import (
	"starbunk-bot/internal/config"
	"starbunk-bot/internal/log"

	"github.com/bwmarrin/discordgo"
)

type ClearWebhooks struct {
	Command string
	GuildID string
}

func (c ClearWebhooks) CommandWord() string {
	return c.Command
}

func (c ClearWebhooks) IsValidCommand(message string) bool {
	return isValidCommand(c.Command, message)
}

func (c ClearWebhooks) ProcessMessage(session *discordgo.Session, message discordgo.Message) {
	if message.Member.Permissions&discordgo.PermissionAdministrator == 0 { // admin command only
		var webhooks, err = session.GuildWebhooks(config.GuildIDs["Starbunk"])
		if err != nil {
			log.ERROR.Println("Error Getting Webhooks", err)
		}
		for _, webhook := range webhooks {
			if webhook.ApplicationID == session.State.User.ID {
				log.WARN.Println("Deleting Webhook ", webhook.Name)
				session.WebhookDelete(webhook.ID)
			}
		}
	}
}
