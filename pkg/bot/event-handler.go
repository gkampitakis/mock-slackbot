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

// https://api.slack.com/events/reaction_added
func reactionAddEvent(slackClient *slack.Client, event *slackevents.ReactionAddedEvent) {
	if event.Reaction != "mock" {
		return
	}

	user, err := slackClient.GetUserInfo(event.ItemUser)
	if err != nil {
		log.Println("[warning]: can't get user info", err)
		return
	}
	if user.IsBot {
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

	message := res.Messages[0].Msg.Text
	_, _, err = slackClient.PostMessage(
		event.Item.Channel,
		slack.MsgOptionBlocks(
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", mock.Mockerize(message), false, true),
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
