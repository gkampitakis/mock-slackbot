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

func healthcheckHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Healthy Status",
	})
}

func eventsEndpointHandler(eventChanel chan<- slackevents.EventsAPIInnerEvent) func(*gin.Context) {
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
