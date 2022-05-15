package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type configuration struct {
	_             struct{}
	SigningSecret string
	BotToken      string
	IsProduction  bool
}

func init() {
	if os.Getenv("BOT_MODE") == "production" {
		return
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Printf("[configuration-init]-[error]: %s\n", err)

		return
	}

	log.Println("[configuration-init]: Loaded")
}

func NewConfiguration() *configuration {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	isProduction := os.Getenv("BOT_MODE") == "production"

	if signingSecret == "" || slackBotToken == "" {
		log.Fatalln("[configuration]: SLACK_SIGNING_SECRET, SLACK_BOT_TOKEN are required env vars")
	}

	return &configuration{
		SigningSecret: signingSecret,
		BotToken:      slackBotToken,
		IsProduction:  isProduction,
	}
}
