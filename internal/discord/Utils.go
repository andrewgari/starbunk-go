package discord

import (
	"net/http"
	"starbunk-bot/internal/log"

	"github.com/bwmarrin/discordgo"
)

type DiscordMemberInfo struct {
	UserID      string
	DisplayName string
	AvatarURL   string
}

func GetMemberInfo(session *discordgo.Session, userID string, guildID string) (DiscordMemberInfo, error) {
	var guildMember, err = session.GuildMember(guildID, userID)
	// if they don't have a nickname, error?
	if err != nil {
		log.WARN.Println("Can't find member info from Server", err)
	}
	return DiscordMemberInfo{UserID: userID, DisplayName: guildMember.Nick, AvatarURL: guildMember.AvatarURL("")}, err
}

func GetFilesFromAttachments(attachments []*discordgo.MessageAttachment) []*discordgo.File {
	var files = make([]*discordgo.File, len(attachments)) // make([]*discordgo.File, len(attachments))
	for _, s := range attachments {
		response, error := http.Get(s.URL)
		if error != nil {
			continue
		}
		log.INFO.Println(s.URL)
		log.INFO.Println(s.ProxyURL)
		var file = &discordgo.File{Name: s.Filename, ContentType: s.ContentType, Reader: response.Body}
		files = append(files, file)
		response.Body.Close()
	}
	return files
}
