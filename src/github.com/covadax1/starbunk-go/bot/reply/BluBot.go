package reply

import (
	"log"
	"regexp"
	"strings"
	"time"

	"golang-discord-bot/src/github.com/covadax1/starbunk-go/bot/config"
	"golang-discord-bot/src/github.com/covadax1/starbunk-go/bot/webhook"

	"github.com/bwmarrin/discordgo"
)

type BluBot struct {
	Name string
}

const defaultPattern string = ".*?\b(blue?|bloo|b lu|eulb|azul|azulbot|cerulean)\b[^$]*$"
const confirmPattern string = ".*?\b(blue?(bot)?)|(bot)|yes|no|yep|(i did)|(you got it)|(sure did)\b[^$]*$"
const nicePattern string = "blue?bot,? say something nice about (?P<name>.+$)"
const meanPattern string = "\b(fuck(ing)?|hate|die|kill|worst|mom|shit|murder|bots?)\b"

const murderAvatar string = "https://imgur.com/Tpo8Ywd.jpg"
const cheekyAvatar string = "https://i.imgur.com/dO4a59n.png"

const defaultResponse string = "Did somebody say Blu?"
const cheekyResponse string = "Lol, Somebody definitely said Blu! :smile:"
const friendlyResponse string = "%s, I think you're pretty Blu! :wink:"
const contemptResponse string = "No way, Venn can suck my blu cane. :unamused:"
const murderResponse string = "What the fuck did you just fucking say about me, you little bitch? I'll have you know I graduated top of my class in the Academia d'Azul, and I've been involved in numerous secret raids on Western La Noscea, and I have over 300 confirmed kills. I've trained with gorillas in warfare and I'm the top bombardier in the entire Eorzean Alliance. You are nothing to me but just another target. I will wipe you the fuck out with precision the likes of which has never been seen before on this Shard, mark my fucking words. You think you can get away with saying that shit to me over the Internet? Think again, fucker. As we speak I am contacting my secret network of tonberries across Eorzea and your IP is being traced right now so you better prepare for the storm, macaroni boy. The storm that wipes out the pathetic little thing you call your life. You're fucking dead, kid. I can be anywhere, anytime, and I can kill you in over seven hundred ways, and that's just with my bear-hands. Not only am I extensively trained in unarmed combat, but I have access to the entire arsenal of the Eorzean Blue Brigade and I will use it to its full extent to wipe your miserable ass off the face of the continent, you little shit. If only you could have known what unholy retribution your little \"clever\" comment was about to bring down upon you, maybe you would have held your fucking tongue. But you couldn't, you didn't, and now you're paying the price, you goddamn idiot. I will fucking cook you like the little macaroni boy you are. You're fucking dead, kiddo."

var bluTimestamp = time.Unix(0, 0)
var bluMurderTimestamp = time.Unix(0, 0)

func (b BluBot) ObserverName() string {
	return "BluBot"
}

func (b BluBot) AvatarURL() string {
	return "https://imgur.com/WcBRCWn.png"
}

func (b BluBot) HandleMessage(session *discordgo.Session, message discordgo.Message) {
	channelID := message.ChannelID
	if strings.Contains(message.Content, "blu") {
		log.Default().Println("Running BlueBot HandleMessage")
		webhook.WriteMessage(session, channelID, "Did somebody say BLU?", b.ObserverName(), b.AvatarURL())
	} else if isRequestToSayBlu(message.Content) {
		getNameFromBluRequest(message.Content, message.Author.Username)
	}
}

func isRequestToSayBlu(message string) bool {
	match, err := regexp.MatchString(nicePattern, message)
	if err != nil {
		log.Fatal("Error Parsing Message: ", err)
	}
	return match
}

func getNameFromBluRequest(message, author string) string {
	regex, err := regexp.Compile(nicePattern)
	if err != nil {
		log.Fatal("Error Parsing Message: ", err)
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
	return message.ReferencedMessage != nil &&
		message.ReferencedMessage.Author.Username == selfID &&
		bluMurderTimestamp.UTC().YearDay() < time.Now().YearDay()
}

func isVennInsultingBlu(message, authorID string) bool {
	if authorID == config.UserIDs["venn"] {
		match, err := regexp.MatchString(meanPattern, message)
		if err != nil {
			log.Fatal("Error Parsing Message: ", err)
		}
		return match
	}
	return false
}
