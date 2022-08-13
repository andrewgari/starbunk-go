package command

import (
	"fmt"
	"starbunk-bot/internal/log"

	"github.com/bwmarrin/discordgo"
)

type MusicCorrect struct {
	Command string
}

func (c MusicCorrect) CommandWord() string {
	return c.Command
}

func (c MusicCorrect) IsValidCommand(message string) bool {
	return isValidCommand(c.Command, message)
}

func (c MusicCorrect) ProcessMessage(session *discordgo.Session, message discordgo.Message) {
	_, err := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Hey <@%s>, Buddy.\n I see you're trying to activate the music bot. I get it, I love to jam it out from time to time. But hey, let me fill you in on a little insider secret. \n Ya see, the bot's gone through some *changes* lately, and some of the functions have changed. What *used* to be `?play` or `?covaPlay` has been updated to just `!play`. If you'd like, you can even ask me directly by saying `@%s play A Fart with Extra Reverb`\n I know! It's that simple, so if you want to jam it out with your buds or just wanna troll them with some cockamamie video of Tidus Laughing to the DK Rap or something (I dunno, I'm not judging) you can call on me anytime with some youtube link", message.Author.ID, session.State.User.Username))
	if err != nil {
		log.ERROR.Println("Error helping a bud out with the music bot")
	}
}
