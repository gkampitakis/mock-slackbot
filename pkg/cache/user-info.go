package cache

import (
	"log"

	"github.com/Code-Hex/go-generics-cache/policy/lfu"
	"github.com/slack-go/slack"
)

var userBotCache = lfu.NewCache[string, bool](lfu.WithCapacity(1024))

func IsUserBot(slackClient *slack.Client, userID string) bool {
	if isBot, exists := userBotCache.Get(userID); exists {
		return isBot
	}

	user, err := slackClient.GetUserInfo(userID)
	if err != nil {
		log.Println("[warning]: can't get user info", err)
		userBotCache.Set(userID, true)
		return true
	}

	userBotCache.Set(userID, user.IsBot)
	return user.IsBot
}
