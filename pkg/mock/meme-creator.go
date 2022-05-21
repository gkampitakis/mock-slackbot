package mock

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	MemeAPIUrl     = "https://api.imgflip.com/caption_image"
	MemeTemplateID = "102723630" // Bob Squarepants mock meme
)

var (
	MemeAPIUsername = os.Getenv("MEME_API_USERNAME")
	MemeAPIPassword = os.Getenv("MEME_API_PASSWORD")
)

type memeResponseObject struct {
	Success      bool              `json:"success"`
	Data         map[string]string `json:"data"`
	ErrorMessage string            `json:"error_message"`
}

func CreateMeme(msg string) (string, error) {
	response, err := http.PostForm(MemeAPIUrl, url.Values{
		"template_id":    {MemeTemplateID},
		"username":       {MemeAPIUsername},
		"password":       {MemeAPIPassword},
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
