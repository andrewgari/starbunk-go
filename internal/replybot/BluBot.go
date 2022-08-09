package replybot

import (
	"fmt"
	"golang-discord-bot/internal/config"
	"golang-discord-bot/internal/log"
	"golang-discord-bot/internal/webhook"
	"golang-discord-bot/utils"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type BluBot struct {
	Name string
}

const (
	defaultPattern string = "(?i)\\b(blue?|bloo|b lue?|eulb|azul|azul|cerulean|azure)(bot)?\\b"
	confirmPattern string = "(?i)\\b(blue?(bot)?)|(bot)|yes|no|yep|yeah|(i did)|(you got it)|(sure did)\\b"
	nicePattern    string = "(?i)blue?bot,? say something nice about (?P<name>.+$)"
	meanPattern    string = "(?i)\\b(fuck(ing)?|hate|die|kill|worst|mom|shit|murder|bots?)\\b"

	murderAvatar string = "https://imgur.com/Tpo8Ywd.jpg"
	cheekyAvatar string = "https://i.imgur.com/dO4a59n.png"

	defaultResponse  string = "Did somebody say Blu?"
	cheekyResponse   string = "Lol, Somebody definitely said Blu! :smile:"
	friendlyResponse string = "%s, I think you're pretty Blu! :wink:"
	contemptResponse string = "No way, Venn can suck my blu cane. :unamused:"
	murderResponse   string = "What the fuck did you just fucking say about me, you little bitch? I'll have you know I graduated top of my class in the Academia d'Azul, and I've been involved in numerous secret raids on Western La Noscea, and I have over 300 confirmed kills. I've trained with gorillas in warfare and I'm the top bombardier in the entire Eorzean Alliance. You are nothing to me but just another target. I will wipe you the fuck out with precision the likes of which has never been seen before on this Shard, mark my fucking words. You think you can get away with saying that shit to me over the Internet? Think again, fucker. As we speak I am contacting my secret network of tonberries across Eorzea and your IP is being traced right now so you better prepare for the storm, macaroni boy. The storm that wipes out the pathetic little thing you call your life. You're fucking dead, kid. I can be anywhere, anytime, and I can kill you in over seven hundred ways, and that's just with my bear-hands. Not only am I extensively trained in unarmed combat, but I have access to the entire arsenal of the Eorzean Blue Brigade and I will use it to its full extent to wipe your miserable ass off the face of the continent, you little shit. If only you could have known what unholy retribution your little \"clever\" comment was about to bring down upon you, maybe you would have held your fucking tongue. But you couldn't, you didn't, and now you're paying the price, you goddamn idiot. I will fucking cook you like the little macaroni boy you are. You're fucking dead, kiddo."
)

var (
	bluTimestamp       = time.Unix(0, 0)
	bluMurderTimestamp = time.Unix(0, 0)
)

func (b BluBot) ObserverName() string {
	return "BluBot"
}

func (b BluBot) AvatarURL() string {
	return "https://imgur.com/WcBRCWn.png"
}

func (b BluBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	channelID := message.ChannelID
	if isRequestToSayBlu(message.Content) {
		name := getNameFromBluRequest(message.Content, message.Author.Username)
		if strings.ToLower(name) == "venn" {
			webhook.WriteMessage(session, channelID, contemptResponse, b.ObserverName(), b.AvatarURL())
		} else {
			webhook.WriteMessage(session, channelID, fmt.Sprintf(friendlyResponse, name), b.ObserverName(), b.AvatarURL())
		}
	} else if isVennInsultingBlu(message.Content, message.Author.ID) {
		bluTimestamp = time.Now()
		bluMurderTimestamp = time.Now()
		webhook.WriteMessage(session, channelID, murderResponse, b.ObserverName(), murderAvatar)
	} else if isResponseToBlu(message, session.State.SessionID) {
		bluTimestamp = time.Unix(0, 0)
		webhook.WriteMessage(session, channelID, cheekyResponse, b.ObserverName(), cheekyAvatar)
	} else if didSomebodySayBlu(message.Content) {
		bluTimestamp = time.Now()
		webhook.WriteMessage(session, channelID, defaultResponse, b.ObserverName(), b.AvatarURL())
	}
}

func didSomebodySayBlu(message string) bool {
	return utils.Match(defaultPattern, message)
}

func confirmSomebodySaidBlu(message string) bool {
	return utils.Match(confirmPattern, message)
}

func isRequestToSayBlu(message string) bool {
	return utils.Match(nicePattern, message)
}

func getNameFromBluRequest(message, author string) string {
	regex, err := regexp.Compile(nicePattern)
	if err != nil {
		log.ERROR.Println("Error Parsing Message: ", err)
	}
	var matches = regex.FindStringSubmatch(message)
	var index = regex.SubexpIndex("name")
	name := "Hey"
	if index > -1 {
		if strings.ToLower(name) == "me" {
			name = author
		} else {
			name = matches[index]
		}
	}
	return name
}

func isResponseToBlu(message discordgo.Message, selfID string) bool {
	if message.ReferencedMessage != nil && message.ReferencedMessage.Author.Username == selfID {
		log.INFO.Println("Message is Referenced by me")
		return true
	} else if message.Timestamp.Before(bluTimestamp.Add(3e+11)) && utils.Match(confirmPattern, message.Content) { // if the message timestamp is within 5 minutes of the last blue message
		return true
	}
	return false
}

func isVennInsultingBlu(message, authorID string) bool {
	if authorID == config.UserIDs["venn"] && bluMurderTimestamp.UTC().Day() < time.Now().Day() && utils.Match(meanPattern, message) {
		return true
	}
	return false
}
