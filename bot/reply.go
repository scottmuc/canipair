package bot

import (
	"fmt"
	"strings"
	"time"
)

// HelpMessage for when responding to help queries
const HelpMessage = `
I am the pairing match maker! I will let you know who in this channel
can pair with you! Here are the things you can ask me:

* today?
* week?
* help!

If you would like to declare your ability to pair, set your Slack emoji status
to anything with the word *beach* in it (e.g.: :beach_with_umbrella: or :beachball:).

You can indicate which specific days you can pair by adding the ICP (I can pair)
bit mask to your status text. Here are some examples:

ICP-1101100 == I can pair Monday, Tuesday, Thursday and Friday
ICP-0011100 == I can pair Wednesday, Thursday, and Friday

*note* if your status emoji is not beach related, then the bitmask will be ignored.

*Known Issues*

* Changing your status doesn't get reflected immediately
* Doesn't work in private channels
* Doesn't respect time zones when factoring in if someone is free
`

func needsTodayReply(text string) bool {
	return strings.Contains(text, "today")
}

func todayReply(users []potentialPair) string {
	userStrings := []string{}
	for _, u := range users {
		if u.onTheBeach {
			if u.isAvailableOn(time.Now().Weekday()) {
				userStrings = append(userStrings, u.displayName)
			}
		}
	}
	if len(userStrings) > 0 {
		return fmt.Sprintf("The following people are free to pair:\n  Today: %s", userStrings)
	}

	return "Unfortunately there's no one available in this channel to pair today."
}

func needsThisWeekReply(text string) bool {
	return strings.Contains(text, "week")
}

func thisWeekReply(users []potentialPair) string {
	return "throw a NotImplementedException here please"
}
