package snowbunk

import "github.com/bwmarrin/discordgo"

type ISnowbunkMessageHandler interface {
	SyncMessage(session *discordgo.Session, message discordgo.Message)
}
