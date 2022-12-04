package mock

import (
	"encoding/json"
	"errors"
	"io"
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

type httpClient interface {
	PostForm(string, url.Values) (*http.Response, error)
}

func CreateMockMeme(client httpClient, msg string) (string, error) {
	response, err := client.PostForm(MemeAPIUrl, url.Values{
		"template_id":    {MemeTemplateID},
		"username":       {MemeAPIUsername},
		"password":       {MemeAPIPassword},
		"boxes[0][text]": {msg},
	})
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var jsonBody memeResponseObject
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		return "", err
	}
	if jsonBody.ErrorMessage != "" {
		return "", errors.New(jsonBody.ErrorMessage)
	}

	return jsonBody.Data["url"], nil
}
