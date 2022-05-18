package bot

import (
	"github.com/gkampitakis/mock-slackbot/pkg/mock"
	"github.com/slack-go/slack"
)

func postMessage(
	slackClient *slack.Client,
	msg, channelID, ts string,
	replyToThread bool,
) error {
	options := []slack.MsgOption{
		slack.MsgOptionBlocks(
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", mock.Mockerize(msg), false, true),
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

	return err
}
