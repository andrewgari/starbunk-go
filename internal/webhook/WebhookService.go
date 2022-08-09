package webhook

import (
	"fmt"
	"golang-discord-bot/internal/log"

	"github.com/bwmarrin/discordgo"
)

func getWebhook(session *discordgo.Session, channelID string) (*discordgo.Webhook, error) {
	var channel, err1 = session.Channel(channelID)
	if err1 != nil {
		log.ERROR.Println("Could not find Channel with ID "+channelID, err1)
	}
	var channelName = fmt.Sprintf("Starbunk-Bot-%s", channel.Name)
	var webhooks, err2 = session.ChannelWebhooks(channelID)
	if err2 != nil {
		log.ERROR.Println("Error getting Webhook for channel "+channelName, err2)
		return nil, err2
	}
	for _, webhook := range webhooks {
		if webhook.ChannelID == channelID && webhook.Name == channelName {
			return webhook, nil
		}
	}
	newWebhook, err3 := session.WebhookCreate(channelID, channelName, "https://pbs.twimg.com/profile_images/421461637325787136/0rxpHzVx_400x400.jpeg")
	if err3 != nil {
		log.ERROR.Println("Error creating Webhook for channel "+channelName, err3)
		return nil, err3
	}
	return newWebhook, nil

}

func WriteMessage(session *discordgo.Session, channelID string, content string, nickname string, avatarURL string) {
	var webhook, err1 = getWebhook(session, channelID)
	if err1 != nil {
		log.ERROR.Println("Error Creating Webhook for channel "+channelID, err1)
	}
	var params = discordgo.WebhookParams{Content: content, Username: nickname, AvatarURL: avatarURL, TTS: false, Files: nil, Components: nil, Embeds: nil, AllowedMentions: nil}
	var wh, errT = session.Webhook(webhook.ID)
	if errT != nil {
		log.ERROR.Println("Error Getting Webhook I just created", errT)
	}
	log.INFO.Println(wh.ID)
	var _, err2 = session.WebhookExecute(webhook.ID, webhook.Token, false, &params)
	if err2 != nil {
		log.ERROR.Println("Error Executing Webhook message: "+content+"\n", err2)
	}
}
