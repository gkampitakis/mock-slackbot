package mock

import (
	"math/rand"
	"strings"
)

func Mockerize(message string) string {
	letters := strings.Split(message, "")
	mockMsg := ""

	for _, letter := range letters {
		random := rand.Intn(100)
		mockMsg += changeCase(letter, random > 50)
	}

	return mockMsg
}

func changeCase(letter string, capital bool) string {
	if capital {
		return strings.ToUpper(letter)
	}

	return strings.ToLower(letter)
}
