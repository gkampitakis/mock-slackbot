package messages

import (
	"log"

	"github.com/slack-go/slack"
)

func Post(
	slackClient *slack.Client,
	msg, channelID, ts string,
	replyToThread bool,
) {
	options := []slack.MsgOption{
		slack.MsgOptionBlocks(
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", msg, false, true),
				nil,
				nil,
			),
		),
	}

	if replyToThread {
		options = append(options, slack.MsgOptionTS(ts))
	}

	_, _, err := slackClient.PostMessage(
		channelID,
		options...,
	)
	if err != nil {
		log.Println("[error]: can't post message", err)
	}
}
