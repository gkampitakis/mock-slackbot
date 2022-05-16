package mock

import (
	"math/rand"
	"regexp"
	"strings"
)

var tagsRegex = regexp.MustCompile("(?i)<[!|@][\\d\\w]+>")

func Mockerize(message string) string {
	letters := strings.Split(escapeTags(message), "")
	mockMsg := ""

	for _, letter := range letters {
		random := rand.Intn(100)
		mockMsg += changeCase(letter, random > 50)
	}

	return mockMsg
}

func escapeTags(message string) string {
	return tagsRegex.ReplaceAllLiteralString(message, "")
}

func changeCase(letter string, capital bool) string {
	if capital {
		return strings.ToUpper(letter)
	}

	return strings.ToLower(letter)
}
