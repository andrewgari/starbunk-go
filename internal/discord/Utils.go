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
		return DiscordMemberInfo{}, err
	}
	return DiscordMemberInfo{UserID: userID, DisplayName: guildMember.Nick, AvatarURL: guildMember.AvatarURL("")}, err
}

func GetFilesFromAttachments(attachments []*discordgo.MessageAttachment) []*discordgo.File {
	files := make([]*discordgo.File, 0) // make([]*discordgo.File, len(attachments))
	for _, s := range attachments {
		response, error := http.Get(s.URL)
		if error != nil {
			continue
		}
		// defer response.Body.Close()
		var file = &discordgo.File{Name: s.Filename, ContentType: s.ContentType, Reader: response.Body}
		files = append(files, file)
	}
	return files
}

func DownloadFileFromUrl(attachment *discordgo.MessageAttachment) (*discordgo.File, error) {
	response, err := http.Get(attachment.ProxyURL)
	if err != nil {
		log.ERROR.Println("Error Getting Response for Image", err)
		return nil, err
	}
	defer response.Body.Close()

	var discordFile = &discordgo.File{
		Name:        attachment.Filename,
		ContentType: attachment.ContentType,
		Reader:      response.Body}

	return discordFile, nil
}

func CheckFileInfo(message discordgo.Message) {
	log.INFO.Println(message)
	for _, ma := range message.Attachments {
		log.INFO.Println("ID", ma.ID)
		log.INFO.Println("Filename", ma.Filename)
		log.INFO.Println("Content Type", ma.ContentType)
		log.INFO.Println("Ephemeral", ma.Ephemeral)
		log.INFO.Println("URL", ma.URL)
		log.INFO.Println("ProxyURL", ma.ProxyURL)
	}
}
