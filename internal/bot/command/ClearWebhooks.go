package command

import (
	"starbunk-bot/internal/log"
	"starbunk-bot/internal/webhook"

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
		channels, err := session.GuildChannels(c.GuildID)
		if err != nil {
			log.ERROR.Println("Error Getting Webhooks", err)
		}
		for _, ch := range channels {
			webhook, error := webhook.GetWebhook(session, ch.ID)
			if error != nil {
				log.ERROR.Println("Error getting Webhook on ", ch.ID, error)
				return
			}
			log.INFO.Println("Webhook.ApplicationID: ", webhook.ApplicationID, ".\t Session.State.SessionID: ", session.State.SessionID)
			if webhook.ApplicationID == session.State.SessionID {
				log.WARN.Println("Deleting Webhook ", webhook.Name)
				session.WebhookDelete(webhook.ID)
			}
		}
	}
}
