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
	messageRef := slack.ItemRef{
		Timestamp: event.TimeStamp,
		Channel:   event.Channel,
	}

	err := slackClient.AddReaction("eyes", messageRef)
	if err != nil {
		log.Println(err)
	}

	mockMsg := mock.Mockerize(event.Text)
	err = postMessage(
		slackClient,
		mockMsg,
		event.Channel,
		event.TimeStamp,
		event.ThreadTimeStamp != "",
	)
	if err != nil {
		log.Println("[error]: can't post message", err)
	}
}

func reactionAddEvent(slackClient *slack.Client, event *slackevents.ReactionAddedEvent) {
	if event.Reaction != "mock" {
		return
	}

	channelID := event.Item.Channel
	ts := event.Item.Timestamp

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

	err = postMessage(
		slackClient,
		// TODO: we can add image here as well ?
		res.Messages[0].Msg.Text,
		channelID,
		ts,
		true,
	)
	if err != nil {
		log.Println("[error]: can't post message", err)
		return
	}
}

func alreadyAnswered(msg slack.Message) bool {
	return msg.ThreadTimestamp != ""
}
