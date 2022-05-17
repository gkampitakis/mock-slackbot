package bot

import (
	"log"

	"github.com/gkampitakis/mock-slackbot/pkg/mock"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func registerHandler(slackClient *slack.Client) func(slackevents.EventsAPIInnerEvent) {
	return func(innerEvent slackevents.EventsAPIInnerEvent) {
		switch event := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			appMentionEvent(slackClient, event)
		case *slackevents.ReactionAddedEvent:
			reactionAddEvent(slackClient, event)
		default:
			log.Println("[warning]: unsupported event")
		}
	}
}

func appMentionEvent(slackClient *slack.Client, event *slackevents.AppMentionEvent) {
	log.Println(event)
}

func reactionAddEvent(slackClient *slack.Client, event *slackevents.ReactionAddedEvent) {
	if event.Reaction != "mock" {
		return
	}

	channelID := event.Item.Channel
	ts := event.Item.Timestamp

	// Not needed
	// if cache.IsUserBot(slackClient, event.ItemUser) {
	// 	return
	// }

	res, err := slackClient.GetConversationHistory(&slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Inclusive: true,
		Latest:    ts,
		Limit:     1,
	})
	if err != nil {
		log.Println("[error]: can't retrieve message", err)
		return
	}

	// The second check is to verify I am answering the message that the emoji was clicked on
	if alreadyAnswered(res.Messages[0]) || ts != res.Messages[0].Timestamp {
		return
	}

	innerText := res.Messages[0].Msg.Text
	_, _, err = slackClient.PostMessage(
		event.Item.Channel,
		slack.MsgOptionBlocks(
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", mock.Mockerize(innerText), false, true),
				nil,
				nil,
			),
		),
		slack.MsgOptionTS(event.Item.Timestamp), // For sending to a thread
	)
	if err != nil {
		log.Println("[error]: can't post message", err)
		return
	}
}

func alreadyAnswered(msg slack.Message) bool {
	return msg.ThreadTimestamp != ""
}
