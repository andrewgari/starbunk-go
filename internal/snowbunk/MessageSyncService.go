package snowbunk

import (
	"starbunk-bot/internal/discord"
	"starbunk-bot/internal/log"
	"starbunk-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
)

var MessageSyncService ISnowbunkMessageHandler

var channelList map[string][]string = map[string][]string{
	"757866614787014660": {"856617421942030364", "798613445301633137"}, // testing
	"856617421942030364": {"757866614787014660", "798613445301633137"}, // testing
	"798613445301633137": {"757866614787014660", "856617421942030364"}, // testing

	"755579237934694420": {"755585038388691127"}, // starbunk
	"755585038388691127": {"755579237934694420"}, // starbunk

	"753251583084724371": {"697341904873979925"}, // memes
	"697341904873979925": {"753251583084724371"}, // memes

	"754485972774944778": {"696906700627640352"}, // ff14 general
	"696906700627640352": {"754485972774944778"}, // ff14 general

	"697342576730177658": {"753251583084724372"}, // ff14 msq
	"753251583084724372": {"697342576730177658"}, // ff14 msq

	"753251583286050926": {"755575759753576498"}, // screenshots
	"755575759753576498": {"753251583286050926"}, // screenshots

	"753251583286050928": {"699048771308224642"}, // raiding
	"699048771308224642": {"753251583286050928"}, // raiding

	"696948268579553360": {"755578695011270707"}, // food
	"755578695011270707": {"696948268579553360"}, // food

	"696948305586028544": {"755578835122126898"}, // pets
	"755578835122126898": {"696948305586028544"}, // pets
}

type SnowbunkService struct {
	StarbunkToken string
	SnowbunkToken string
}

func (snowservice SnowbunkService) SyncMessage(session *discordgo.Session, message discordgo.Message) {
	var originChannel, err = session.Channel(message.ChannelID)
	if err != nil {
		log.ERROR.Println("Can't find Channen Origin", err)
		return
	}

	for _, link := range channelList[message.ChannelID] {
		var linkedChannel, err2 = session.Channel(link)
		if err2 != nil {
			log.ERROR.Println("Can't find Linked Channel", err2)
			continue
		}
		// write to webhook on channel
		snowservice.WriteMessage(session, message, originChannel, linkedChannel)
	}
}

func (snowservice SnowbunkService) WriteMessage(session *discordgo.Session, message discordgo.Message, originChannel *discordgo.Channel, linkedChannel *discordgo.Channel) {
	var info, err = discord.GetMemberInfo(session, message.Author.ID, linkedChannel.GuildID)
	if err != nil || info.DisplayName == "" {
		log.WARN.Println("Can't find Member Info at Linked Server")
		info, err = discord.GetMemberInfo(session, message.Author.ID, originChannel.GuildID)
		if err != nil || info.DisplayName == "" {
			log.WARN.Println("Can't Find User At Origin Server")
			info = discord.DiscordMemberInfo{UserID: message.Author.ID, DisplayName: message.Author.Username, AvatarURL: message.Author.AvatarURL("")}
		}
	}
	webhook.WriteMessage(session, linkedChannel.ID, message.Content, info.DisplayName, info.AvatarURL, message.Attachments)
}
