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
	return "Hey <@%s>, Buddy.\nI see you're trying to activate the music bot... I get it, I love to jam it out from time to time. But hey, let me fill you in on a little insider secret.\nYa see, the bot's gone through even **more** *changes* lately (Yeah, Yeah, I know. It keeps on changing how can my tiny brain keep up :unamused:). What *used* to be `?play` or `!play` has been updated to the shiny new command `/play`.\nIf you'd like, you can even ask me directly by saying `@%s play A Fart with Extra Reverb`\nI know! It's that simple, so if you want to jam it out with your buds or just wanna troll them with some stupid video of a gross man in dirty underpants farting on his roomate's door or .... just the sound of a fart with a little extra revery (I dunno, I'm not judging :shrug:) you can call on me anytime with some youtube link.\n"
}

func (c MusicCorrect) Pattern() string {
	return "(\\?|\\!)play "
}

func (c MusicCorrect) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if utils.Match(c.Pattern(), message.Content) {
		_, err := session.ChannelMessageSend(
			message.ChannelID,
			fmt.Sprintf(c.Response(), message.Author.ID, session.State.User.Username))
		if err != nil {
			log.ERROR.Println("Error helping a bud out with the music bot", err)
		}
	}
}
