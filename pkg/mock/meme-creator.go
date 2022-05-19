package mock

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	MEME_API_URL     = "https://api.imgflip.com/caption_image"
	MEME_TEMPLATE_ID = "102723630" // Bob Squarepants mock meme
)

var (
	MEME_API_USERNAME = os.Getenv("MEME_API_USERNAME")
	MEME_API_PASSWORD = os.Getenv("MEME_API_PASSWORD")
)

type memeResponseObject struct {
	Success      bool              `json:"success"`
	Data         map[string]string `json:"data"`
	ErrorMessage string            `json:"error_message"`
}

func CreateMeme(msg string) (string, error) {
	response, err := http.PostForm(MEME_API_URL, url.Values{
		"template_id":    {MEME_TEMPLATE_ID},
		"username":       {MEME_API_USERNAME},
		"password":       {MEME_API_PASSWORD},
		"boxes[0][text]": {msg},
	})
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var jsonBody memeResponseObject
	err = json.Unmarshal(body, &jsonBody)
	if err != nil || jsonBody.ErrorMessage != "" {
		return "", err
	}

	return jsonBody.Data["url"], nil
}
