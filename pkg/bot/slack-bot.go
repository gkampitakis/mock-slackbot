package bot

import (
	"github.com/gin-gonic/gin"
	"github.com/gkampitakis/mock-slackbot/pkg/utils"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var config = utils.NewConfiguration()

type Bot struct {
	_            struct{}
	slackClient  *slack.Client
	api          *gin.Engine
	eventChannel chan slackevents.EventsAPIInnerEvent
}

func (bot *Bot) registerRoutes() {
	bot.api.GET("/health", healthcheckHandler)
	bot.api.POST("/events-endpoint", eventsEndpointHandler(bot.eventChannel))
}

func (bot *Bot) Run() error {
	return bot.api.Run()
}

func (bot *Bot) eventLoop() {
	handler := registerHandler(bot.slackClient)

	for event := range bot.eventChannel {
		handler(event)
	}
}

// TODO: graceful shutdown ????

func NewBot() *Bot {
	bot := &Bot{
		slackClient:  slack.New(config.BotToken),
		api:          gin.Default(),
		eventChannel: make(chan slackevents.EventsAPIInnerEvent, 1024),
	}

	bot.registerRoutes()
	// Running on it's own go routine
	go bot.eventLoop()

	if config.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	return bot
}
