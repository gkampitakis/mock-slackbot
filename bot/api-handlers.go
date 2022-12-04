package bot

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var uptime = time.Now()

func healthcheckHandler(ctx *gin.Context) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	ctx.JSON(200, gin.H{
		"service":     "mock-bot",
		"uptime":      time.Since(uptime).String(),
		"go-routines": runtime.NumGoroutine(),
		"memory": map[string]interface{}{
			"rss":                mem.HeapSys,
			"total-alloc":        mem.TotalAlloc,
			"heap-alloc":         mem.HeapAlloc,
			"heap-objects-count": mem.HeapObjects,
		},
	})
}

func eventsEndpointHandler(eventChanel chan<- slackevents.EventsAPIInnerEvent) func(*gin.Context) {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
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

		/**
			De-couple your ingestion of events from processing and reacting to them.
			Especially when working with large workspaces, many workspaces, or subscribing to a large number of events.
			Quickly respond to events with HTTP 200 and add them to a queue before doing amazing things with them.
		**/
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

	log.Println("server verified")
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
