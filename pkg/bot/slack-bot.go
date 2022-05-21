package bot

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	server       *http.Server
}

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
		bot.gracefulShutdown()
		gin.SetMode(gin.ReleaseMode)

		// Running on it's own go routine
		go bot.init()
	}

	return bot
}

// TODO: check params are not used go lint
func (bot *Bot) init() {
	log.Println("[info]: initializing slack-bot")

	channels, _, err := bot.slackClient.GetConversations(
		&slack.GetConversationsParameters{
			ExcludeArchived: true,
			Types:           []string{"public_channel"},
		})
	if err != nil {
		log.Printf("[error]: can't retrieve channels %s\n", err)
		return
	}

	utils.Concurrent(channels, func(channel slack.Channel) {
		_, _, _, err := bot.slackClient.JoinConversation(channel.ID)
		if err != nil {
			log.Printf("[error]: can't join channel %s, %s\n", channel.Name, err)
		}
	}, 20)
}

func (bot *Bot) registerRoutes() {
	bot.api.GET("/health", healthcheckHandler)
	bot.api.POST("/events-endpoint", eventsEndpointHandler(bot.eventChannel))
}

func (bot *Bot) Run() error {
	bot.server = &http.Server{
		Addr:    ":8080",
		Handler: bot.api,
	}

	log.Println("Listening and serving HTTP on :8080")

	return bot.server.ListenAndServe()
}

func (bot *Bot) eventLoop() {
	parallel := make(chan struct{}, 50)
	handler := registerHandler(bot.slackClient)

	for event := range bot.eventChannel {
		parallel <- struct{}{}

		go func(e slackevents.EventsAPIInnerEvent) {
			handler(e)
			<-parallel
		}(event)
	}
}

func (bot *Bot) gracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(
		c,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {
		<-c

		time.Sleep(1 * time.Second)

		for {
			if len(bot.eventChannel) == 0 {
				break
			}

			log.Println("draining event queue")
			time.Sleep(1 * time.Second)
		}

		log.Println("[info]: bot shutting down")

		if err := bot.server.Shutdown(context.TODO()); err != nil {
			log.Fatalln(err)
		}
	}()
}
