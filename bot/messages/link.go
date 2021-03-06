package messages

import (
	"log"

	"github.com/slack-go/slack"
)

func Link(
	slackClient SlackClient,
	link, channelID, ts string,
	replyToThread bool,
) {
	options := []slack.MsgOption{
		slack.MsgOptionText(link, false),
	}

	if replyToThread {
		options = append(options, slack.MsgOptionTS(ts))
	}

	_, _, err := slackClient.PostMessage(
		channelID,
		options...,
	)
	if err != nil {
		log.Println("[error]: can't post link message", err)
	}
}
