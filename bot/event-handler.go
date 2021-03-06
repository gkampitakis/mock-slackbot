package bot

import (
	"log"
	"strings"

	"github.com/gkampitakis/mock-slackbot/bot/messages"
	"github.com/gkampitakis/mock-slackbot/mock"
	"github.com/gkampitakis/mock-slackbot/utils"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	MuteMeme = "https://imgflip.com/i/6gn6ir"
)

var users = map[string]struct{}{}

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

	msg := utils.EscapeSlackTags(event.Text)
	if msg == "" {
		log.Println("[warning]: not txt message")
		return
	}

	if _, exists := users[event.User]; !exists && isMuteCommand(msg) {
		messages.Link(
			slackClient,
			MuteMeme,
			event.Channel,
			event.TimeStamp,
			event.ThreadTimeStamp != "",
		)

		users[event.User] = struct{}{}
		return
	}

	messages.Post(
		slackClient,
		mock.Mockerize(msg),
		event.Channel,
		event.TimeStamp,
		event.ThreadTimeStamp != "",
	)
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

	msg := utils.EscapeSlackTags(res.Messages[0].Msg.Text)
	if msg == "" {
		log.Println("[warning]: not txt message")
		return
	}

	messages.Post(
		slackClient,
		// TODO: we can add image here as well ?
		mock.Mockerize(msg),
		channelID,
		ts,
		true,
	)

	if _, exists := users[event.ItemUser]; exists {
		return
	}

	messages.Mute(slackClient, channelID, event.ItemUser, ts)
}

func alreadyAnswered(msg slack.Message) bool {
	return msg.ThreadTimestamp != ""
}

func isMuteCommand(msg string) bool {
	words := strings.Fields(msg)

	if len(words) == 1 && words[0] == "mute" {
		return true
	}

	return false
}
