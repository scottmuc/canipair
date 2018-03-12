package main

import (
	"fmt"
	"os"

	"github.com/scottmuc/canipair/bot"
)

func main() {
	token, isSet := os.LookupEnv("SLACK_API_TOKEN")
	if !isSet {
		fmt.Println("Please set SLACK_API_TOKEN environment variable")
		os.Exit(1)
	}

	bot.New(token).Run()
}
