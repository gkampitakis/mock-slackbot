package mock

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
)

func TestMockerize(t *testing.T) {
	t.Run("should randomly change word casing", func(t *testing.T) {
		mockMsg := "some mock message"
		result := Mockerize(mockMsg)

		assert.NotEqual(t, result, mockMsg)
		assert.Equal(t, strings.ToLower(result), strings.ToLower(mockMsg))
	})
}

type MockPostFrom func(string, url.Values) (*http.Response, error)

type mockClient struct {
	mockPostForm MockPostFrom
}

func (c mockClient) PostForm(url string, values url.Values) (*http.Response, error) {
	return c.mockPostForm(url, values)
}

func getMockClient(mockFn MockPostFrom) mockClient {
	return mockClient{
		mockPostForm: mockFn,
	}
}

func getMockResponse(response string) *http.Response {
	return &http.Response{
		Body: io.NopCloser(strings.NewReader(response)),
	}
}

func TestCreateMockMeme(t *testing.T) {
	t.Run("should call postFrom with correct values", func(t *testing.T) {
		mockURL := "mock-url"
		mock := func(url string, values url.Values) (*http.Response, error) {
			snaps.MatchSnapshot(t, url, values)

			return getMockResponse(
				fmt.Sprintf(`{"success":true,"data":{"url":"%s"}}`, mockURL),
			), nil
		}

		url, err := CreateMockMeme(getMockClient(mock), "mock-message")

		assert.Nil(t, err)
		assert.Equal(t, mockURL, url)
	})

	t.Run("should return PostForm error", func(t *testing.T) {
		mockError := errors.New("mock-error")
		mock := func(url string, values url.Values) (*http.Response, error) {
			snaps.MatchSnapshot(t, url, values)

			return nil, mockError
		}

		url, err := CreateMockMeme(getMockClient(mock), "mock-message")

		assert.ErrorIs(t, mockError, err)
		assert.Empty(t, url)
	})

	t.Run("should return error from api", func(t *testing.T) {
		mock := func(url string, values url.Values) (*http.Response, error) {
			snaps.MatchSnapshot(t, url, values)

			return getMockResponse(`{"success":false,"error_message":"mock-error"}`), nil
		}

		url, err := CreateMockMeme(getMockClient(mock), "mock-message")

		assert.Empty(t, url)
		assert.Equal(t, "mock-error", err.Error())
	})
}
