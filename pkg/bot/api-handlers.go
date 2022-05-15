package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// FIXME: print server details
func healthcheckHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Healthy Status",
	})
}

// comment details
// TODO: here we want a channel to post events
// De-couple your ingestion of events from processing and reacting to them.
// Especially when working with large workspaces, many workspaces, or subscribing to a large number of events.
// Quickly respond to events with HTTP 200 and add them to a queue before doing amazing things with them.

func eventsEndpointHandler(eventChanel chan<- interface{}) func(*gin.Context) {
	return func(ctx *gin.Context) {
		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			badRequest(ctx, err)
			return
		}

		if !verifyRequest(ctx, body) {
			return
		}

		slackEvent, err := slackevents.ParseEvent(
			json.RawMessage(body),
			slackevents.OptionNoVerifyToken(),
		)
		if err != nil {
			internalServerError(ctx, err)
			return
		}

		if slackEvent.Type == slackevents.URLVerification {
			slackChallenge(ctx, body)
			return
		}

		if len(eventChanel) == cap(eventChanel) {
			log.Println("[warning]: eventChannel is full! Blocking")
		}
		eventChanel <- slackEvent.InnerEvent

		ctx.Writer.WriteHeader(200)
	}
}

func verifyRequest(ctx *gin.Context, body []byte) bool {
	sv, err := slack.NewSecretsVerifier(ctx.Request.Header, config.SigningSecret)
	if err != nil {
		badRequest(ctx, err)
		return false
	}

	if _, err := sv.Write(body); err != nil {
		internalServerError(ctx, err)
		return false
	}

	if err := sv.Ensure(); err != nil {
		unauthorized(ctx, err)
		return false
	}

	return true
}

func slackChallenge(ctx *gin.Context, body []byte) {
	var r *slackevents.ChallengeResponse

	err := json.Unmarshal(body, &r)
	if err != nil {
		internalServerError(ctx, err)
		return
	}

	ctx.Data(200, "Text", []byte(r.Challenge))
}

func internalServerError(ctx *gin.Context, err error) {
	log.Println(err)

	ctx.JSON(http.StatusInternalServerError, "Internal server error")
}

func badRequest(ctx *gin.Context, err error) {
	log.Println(err)

	ctx.JSON(http.StatusBadRequest, "Bad Request")
}

func unauthorized(ctx *gin.Context, err error) {
	log.Println(err)

	ctx.JSON(http.StatusUnauthorized, "Unauthorized")
}

// get into channels

/*
	https://api.slack.com/methods/conversations.list
	Bot Tokens

	channels:read
	groups:read
	im:read
	mpim:read

	Params

	exclude_archived: boolean
	types: "public_channel"

*/

/**

	channels, _, err := slackAPI.GetConversations(&slack.GetConversationsParameters{
		Types:           []string{"public_channel"},
		ExcludeArchived: true,
	})


app.POST("events-endpoint", func(ctx *gin.Context) {


		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				slackAPI.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			}
		}
	})

*/
