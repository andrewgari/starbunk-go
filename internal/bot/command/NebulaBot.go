package command

import (
	"fmt"
	"starbunk-bot/internal/log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type NebulaBot struct {
	Command        string
	NebulaLeadRole string
	AllowedRoles   map[string]string
}

func (c NebulaBot) CommandWord() string {
	return c.Command
}

func (c NebulaBot) IsValidCommand(message string) bool {
	return isValidCommand(c.Command, message)
}

func (c NebulaBot) HasPermissions(message discordgo.Message) bool {
	return c.contains(message.Member.Roles, c.NebulaLeadRole)
}

func (c NebulaBot) ProcessMessage(session *discordgo.Session, message discordgo.Message) {
	msg := strings.Split(strings.TrimSpace(message.Content), " ")

	if len(msg) < 2 {
		log.ERROR.Printf("Venn made a bad command: %s", message.Content)
		return
	}

	command := msg[1]

	if command == "help" {
		// Help Venn
		c.sendHelpMessage(session, message.ChannelID)
		return
	}

	if len(msg) < 4 {
		log.ERROR.Printf("Venn made a bad command: %s", message.Content)
		return
	}

	guildID := message.GuildID
	memberID := msg[2]
	roleID := msg[3]

	memberID = strings.TrimPrefix(memberID, "<@")
	memberID = strings.TrimSuffix(memberID, ">")
	roleID = strings.TrimPrefix(roleID, "<@&")
	roleID = strings.TrimSuffix(roleID, ">")

	if !c.checkForValue(roleID, c.AllowedRoles) {
		log.ERROR.Printf("Venn going mad with power: %s, %s, %s", message.Content, memberID, roleID)
		return
	}

	switch command {
	case "add":
		session.GuildMemberRoleAdd(guildID, memberID, roleID)
		c.sendFunctionNotice(session, true, message.ChannelID, message.Author.ID, memberID, roleID)
		break
	case "remove":
		session.GuildMemberRoleRemove(guildID, memberID, roleID)
		c.sendFunctionNotice(session, false, message.ChannelID, message.Author.Username, memberID, roleID)
		break
	default:
		log.WARN.Printf("Venn made a bad command: %s", message.Content)
	}

}

func (c NebulaBot) contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// function to check if a value present in the map
func (c NebulaBot) checkForValue(val string, ids map[string]string) bool {
	for _, value := range ids {
		if value == val {
			return true
		}
	}
	return false
}

func (c NebulaBot) sendHelpMessage(session *discordgo.Session, channelID string) {
	_, err := session.ChannelMessageSend(
		channelID,
		"Available Roles: `@Nebula`, `@NebulaFriend`, `@Gobule`\nGrant Role: `!nebula add @<user> @<role>`\nRemove Role: `!nebula remove @<user> @<role>`")
	if err != nil {
		log.ERROR.Println("Error adding Nebula Roles")
	}
}

func (c NebulaBot) sendFunctionNotice(session *discordgo.Session, add bool, channelID, callerID, memberID, roleID string) {
	var function string = ""
	if add {
		function = "added"
	} else {
		function = "removed"
	}

	if function == "added" {
		_, err := session.ChannelMessageSend(
			channelID,
			fmt.Sprintf("<@%s> has added the <@&%s> role to <@%s>", callerID, roleID, memberID))
		if err != nil {
			log.ERROR.Println("Error adding Nebula Roles")
		}
	}
}
