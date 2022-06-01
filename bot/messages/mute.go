package messages

import (
	"log"

	"github.com/slack-go/slack"
)

func Mute(
	slackClient *slack.Client,
	channelID, userID, ts string,
) {
	_, err := slackClient.PostEphemeral(
		channelID,
		userID,
		slack.MsgOptionBlocks(
			slack.NewSectionBlock(
				slack.NewTextBlockObject(
					"mrkdwn",
					"> You can mute me by running @mock-bot mute",
					false,
					true,
				),
				nil,
				nil,
			),
		),
		slack.MsgOptionTS(ts),
	)

	if err != nil {
		log.Println("[error]: can't post ephemeral message", err)
	}
}
