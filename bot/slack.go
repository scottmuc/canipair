package bot

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

const availabilityBitmask = `ICP-(?P<Mon>[01])(?P<Tue>[01])(?P<Wed>[01])(?P<Thu>[01])(?P<Fri>[01])(?P<Sat>[01])(?P<Sun>[01])`

// Runner starts running the initialized bot
type Runner interface {
	Run()
}

type bot struct {
	client *slack.Client
	rtm    *slack.RTM
}

type potentialPair struct {
	displayName   string
	onTheBeach    bool
	availableDays []time.Weekday
}

func (p potentialPair) isAvailableOn(d time.Weekday) bool {
	for _, pd := range p.availableDays {
		if pd == d {
			return true
		}
	}
	return false
}

// New is the bot constructor
func New(token string) Runner {
	api := slack.New(token)
	rtm := api.NewRTM()
	return &bot{
		client: api,
		rtm:    rtm,
	}
}

func handleConnectedEvent(ev *slack.ConnectedEvent) {
	fmt.Println("Connection counter:", ev.ConnectionCount)
}

func handleRTMError(ev *slack.RTMError) {
	fmt.Printf("Error: %s\n", ev.Error())
}

func handleIvalidAuthEvent(ev *slack.InvalidAuthEvent) {
	fmt.Printf("Invalid credentials")
}

func ignoreMessage(reqUser string, botUser string, txt string) bool {
	if reqUser == botUser {
		return true
	}

	prefix := fmt.Sprintf("<@%s>", botUser)
	if !strings.HasPrefix(txt, prefix) {
		return true
	}

	return false
}

func parseAvailableWeekdays(statusText string) []time.Weekday {
	weekDays := []time.Weekday{}
	r := regexp.MustCompile(availabilityBitmask)
	pairingMask := r.FindStringSubmatch(statusText)
	if pairingMask != nil {
		if pairingMask[1] == "1" {
			weekDays = append(weekDays, time.Monday)
		}
		if pairingMask[2] == "1" {
			weekDays = append(weekDays, time.Tuesday)
		}
		if pairingMask[3] == "1" {
			weekDays = append(weekDays, time.Wednesday)
		}
		if pairingMask[4] == "1" {
			weekDays = append(weekDays, time.Thursday)
		}
		if pairingMask[5] == "1" {
			weekDays = append(weekDays, time.Friday)
		}
		if pairingMask[6] == "1" {
			weekDays = append(weekDays, time.Saturday)
		}
		if pairingMask[7] == "1" {
			weekDays = append(weekDays, time.Sunday)
		}
	} else {
		weekDays = []time.Weekday{
			time.Monday,
			time.Tuesday,
			time.Wednesday,
			time.Thursday,
			time.Friday,
			time.Saturday,
			time.Sunday,
		}
	}
	return weekDays
}

func handleMessageEvent(b *bot, ev *slack.MessageEvent) {
	fmt.Printf("Message: %v\n", ev)
	info := b.rtm.GetInfo()
	botUserID := info.User.ID

	if ignoreMessage(ev.User, botUserID, ev.Text) {
		return
	}

	channelInfo, err := b.client.GetChannelInfo(ev.Channel)
	if err != nil {
		fmt.Println(err)
	}

	users := []potentialPair{}

	for _, m := range channelInfo.Members {
		u := info.GetUserByID(m)
		pu := potentialPair{
			displayName:   u.Name,
			onTheBeach:    strings.Contains(u.Profile.StatusEmoji, "beach"),
			availableDays: parseAvailableWeekdays(u.Profile.StatusText),
		}
		users = append(users, pu)
	}

	reply := ""
	if needsTodayReply(ev.Text) {
		reply = todayReply(users)
	} else if needsThisWeekReply(ev.Text) {
		reply = thisWeekReply(users)
	} else {
		reply = HelpMessage
	}
	fmt.Println("Sending reply:", reply)
	b.rtm.SendMessage(b.rtm.NewOutgoingMessage(reply, ev.Channel))
}

func (b *bot) Run() {
	go b.rtm.ManageConnection()

	for {
		select {
		case msg := <-b.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				handleConnectedEvent(ev)
			case *slack.MessageEvent:
				handleMessageEvent(b, ev)
			case *slack.RTMError:
				handleRTMError(ev)
			case *slack.InvalidAuthEvent:
				handleIvalidAuthEvent(ev)
			default:
				//Take no action
			}
		}
	}
}
