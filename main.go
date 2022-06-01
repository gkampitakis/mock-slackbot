package main

import (
	"log"

	"github.com/gkampitakis/mock-slackbot/bot"
)

func main() {
	if err := bot.NewBot().Run(); err != nil {
		log.Fatalln(err)
	}
}
