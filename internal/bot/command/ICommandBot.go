package command

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ICommandBot interface {
	CommandWord() string
	IsValidCommand(message string) bool
	ProcessMessage(session *discordgo.Session, message discordgo.Message)
}

const CommandCharacter string = "!"

func isValidCommand(cmdWord string, message string) bool {
	return strings.HasPrefix(message, CommandCharacter+cmdWord)
}
