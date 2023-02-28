package command

import (
	"fmt"
	"starbunk-bot/internal/log"
	"starbunk-bot/internal/utils"

	"github.com/bwmarrin/discordgo"
)

type MusicCorrect struct{}

func (c MusicCorrect) ObserverName() string {
	return ""
}

func (c MusicCorrect) AvatarURL() string {
	return ""
}

func (c MusicCorrect) Response() string {
	return "Hey <@%s>, Buddy.\nI see you're trying to activate the music bot... I get it, I love to jam it out from time to time. But hey, let me fill you in on a little insider secret.\nYa see, the bot's gone through some *changes* lately (actually it's been like this for a really long time), and some of the functions have changed. What *used* to be `?play` or `?covaPlay` has been updated to just `!play`.\nIf you'd like, you can even ask me directly by saying `@%s play A Fart with Extra Reverb`\nI know! It's that simple, so if you want to jam it out with your buds or just wanna troll them with some cockamamie video of Tidus Laughing to the DK Rap or something (I dunno, I'm not judging) you can call on me anytime with some youtube link.\n"
}

func (c MusicCorrect) Pattern() string {
	return "\\?play "
}

func (c MusicCorrect) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if utils.Match(c.Pattern(), message.Content) {
		_, err := session.ChannelMessageSend(
			message.ChannelID,
			fmt.Sprintf(c.Response(), message.Author.ID, session.State.User.Username))
		if err != nil {
			log.ERROR.Println("Error helping a bud out with the music bot")
		}
	}
}
