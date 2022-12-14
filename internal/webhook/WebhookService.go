package webhook

import (
	"fmt"
	"starbunk-bot/internal/discord"
	"starbunk-bot/internal/log"

	"github.com/bwmarrin/discordgo"
)

func GetWebhook(session *discordgo.Session, channelID, token string) (*discordgo.Webhook, error) {
	var channel, err1 = session.Channel(channelID)
	if err1 != nil {
		log.ERROR.Println("Could not find Channel with ID "+channelID, err1)
	}
	var channelName = fmt.Sprintf("Starbunk-Bot-%s", channel.Name)

	var webhooks, err2 = session.ChannelWebhooks(channelID)
	if err2 != nil {
		log.ERROR.Println("Error getting Webhooks for channel "+channelName, err2)
		return nil, err2
	}
	for _, webhook := range webhooks {
		if webhook.ChannelID == channelID && webhook.Name == channelName {
			return webhook, nil
		}
	}

	log.INFO.Println("Could not find Webhook, Creating New One")
	newWebhook, err3 := session.WebhookCreate(channelID, channelName, "https://pbs.twimg.com/profile_images/421461637325787136/0rxpHzVx_400x400.jpeg")
	if err3 != nil {
		log.ERROR.Println("Error creating Webhook for channel "+channelName, err3)
		return nil, err3
	}

	return newWebhook, nil

}

func WriteMessage(session *discordgo.Session, token string, channelID, content, nickname, avatarURL string, attachments []*discordgo.MessageAttachment) {
	var webhook, err1 = GetWebhook(session, channelID, token)
	if err1 != nil {
		log.ERROR.Println("Error Creating Webhook for channel "+channelID, err1)
	}

	var files = discord.GetFilesFromAttachments(attachments)
	if len(files) > 0 {
		log.INFO.Println("Got Files!")
	} else {
		files = nil
	}
	var params = discordgo.WebhookParams{Content: content, Username: nickname, AvatarURL: avatarURL, TTS: false, Files: files, Components: nil, Embeds: nil, AllowedMentions: nil}
	if token == "" {
		log.WARN.Println("Could not find Webhook Token!")
		log.WARN.Println(webhook)
	}

	var _, err2 = session.WebhookExecute(webhook.ID, webhook.Token, true, &params)
	if err2 != nil {
		log.ERROR.Println("Error Executing Webhook message: "+content+"\n", err2)
	}
}
