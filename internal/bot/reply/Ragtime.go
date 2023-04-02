package reply

import (
	"starbunk-bot/internal/log"
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
	return "http://ragtime.jpeg"
}

func (b RagtimeBot) Response() string {
	return "I am bot"
}

func (b RagtimeBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	channel, err := session.Channel(b.TriviaReviewChannel)
	if err != nil {
		log.ERROR.Println("Couldn't find channel: " + b.TriviaReviewChannel)
		log.WARN.Println("Looking For Channel: " + b.TriviaReviewChannel)
	}
	log.INFO.Println("Message received from Channel: " + channel.Name)
	if message.ChannelID == b.TriviaChannel {
		log.INFO.Println("I'm doin it")
		memberNick := message.Member.Nick
		if memberNick == "" {
			memberNick = message.Author.Username
		}
		if contains(message.Member.Roles, b.TriviaMaster) {
			log.INFO.Println("Trivia Master")
			session.ChannelMessageSend(b.TriviaReviewChannel, "__**"+message.Content+"**__")
			session.ChannelMessageSend(b.TriviaChannel, "You have 20 seconds to answer.")
			time.Sleep(10 * time.Second)
			session.ChannelMessageSend(b.TriviaChannel, "You have 10 seconds left!")
			time.Sleep(10 * time.Second)
			session.ChannelMessageSend(b.TriviaChannel, "*Time's up!*")
			session.ChannelMessageSend(b.TriviaReviewChannel, "*Time's up!*")
		} else {
			log.INFO.Println("Not the Trivia Master")
			session.ChannelMessageDelete(message.ChannelID, message.ID)
			session.ChannelMessageSend(b.TriviaReviewChannel, "**"+memberNick+":** "+message.Content)
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
