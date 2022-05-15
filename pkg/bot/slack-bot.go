package bot

import (
	"github.com/gin-gonic/gin"
	"github.com/gkampitakis/mock-slackbot/pkg/utils"
	"github.com/slack-go/slack"
)

var config = utils.NewConfiguration()

type Bot struct {
	_            struct{}
	slackClient  *slack.Client
	api          *gin.Engine
	eventChannel chan interface{}
}

func (bot *Bot) registerRoutes() {
	bot.api.GET("/health", healthcheckHandler)
	bot.api.POST("/events-endpoint", eventsEndpointHandler(bot.eventChannel))
}

func (bot *Bot) Run() error {
	// FIXME: add some message here :thinking:
	return bot.api.Run()
}

func (bot *Bot) eventLoop() {
	for event := range bot.eventChannel {
		eventHandler(event)
	}
}

// TODO: graceful shutdown ????

func NewBot() *Bot {
	bot := &Bot{
		slackClient:  slack.New(config.BotToken),
		api:          gin.Default(),
		eventChannel: make(chan interface{}, 1024),
	}

	bot.registerRoutes()
	// Running on it's own go routine
	go bot.eventLoop()

	if config.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	return bot
}
