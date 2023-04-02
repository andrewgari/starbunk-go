package reply

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type RagtimeBot struct {
	TriviaMaster        string
	TriviaChannel       string
	TriviaReviewChannel string
}

func (b RagtimeBot) ObserverName() string {
	return "RagtimeBot"
}

func (b RagtimeBot) AvatarURL() string {
	return ""
}

func (b RagtimeBot) Response() string {
	return ""
}

func (b RagtimeBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	if message.ChannelID == b.TriviaChannel {
		memberNick := getNickname(message)
		responses := make(map[string]string)
		if contains(message.Member.Roles, b.TriviaMaster) {
			session.ChannelMessageSend(b.TriviaReviewChannel, fmt.Sprintf("__**%s**__", message.Content))
			session.ChannelMessageSend(b.TriviaChannel, "You have 20 seconds to answer.")
			time.Sleep(10 * time.Second)
			session.ChannelMessageSend(b.TriviaChannel, "You have 10 seconds left!")
			time.Sleep(10 * time.Second)
			session.ChannelMessageSend(b.TriviaChannel, "*Time's up!*")
			session.ChannelMessageSend(b.TriviaReviewChannel, "*Time's up!*")
			for key, value := range responses {
				session.ChannelMessageSend(b.TriviaChannel, fmt.Sprintf("**%s**: *%s*", key, value))
			}
		} else {
			responses[memberNick] = message.Content
			session.ChannelMessageDelete(message.ChannelID, message.ID)
			session.ChannelMessageSend(b.TriviaReviewChannel, fmt.Sprintf("**%s:** *%s*", memberNick, message.Content))
		}
	}
}

func contains(slice []string, check string) bool {
	for _, element := range slice {
		if element == check {
			return true
		}
	}
	return false
}

func getNickname(message discordgo.Message) string {
	nick := message.Member.Nick
	if nick == "" {
		nick = message.Author.Username
	}
	return nick
}
