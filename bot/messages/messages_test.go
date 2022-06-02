package messages

import (
	"testing"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	postMessage   func(string, ...slack.MsgOption) (string, string, error)
	postEphemeral func(string, string, ...slack.MsgOption) (string, error)
}

func (c mockClient) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	return c.postMessage(channelID, options...)
}

func (c mockClient) PostEphemeral(channelID, userID string, options ...slack.MsgOption) (string, error) {
	return c.postEphemeral(channelID, userID, options...)
}

const (
	mockChannelID = "mock-channel"
	mockTS        = "mock-ts"
)

func TestLink(t *testing.T) {
	mockLink := "mock-link"

	t.Run("should pass arguments correctly", func(t *testing.T) {
		c := mockClient{
			postMessage: func(channelID string, options ...slack.MsgOption) (string, string, error) {
				assert.Len(t, options, 1)
				assert.Equal(t, mockChannelID, channelID)

				return "", "", nil
			},
		}

		Link(c, mockLink, mockChannelID, mockTS, false)
	})

	t.Run("should apply replyToThread option", func(t *testing.T) {
		c := mockClient{
			postMessage: func(channelID string, options ...slack.MsgOption) (string, string, error) {
				assert.Len(t, options, 2)
				assert.Equal(t, mockChannelID, channelID)

				return "", "", nil
			},
		}

		Link(c, mockLink, mockChannelID, mockTS, true)
	})
}

func TestMute(t *testing.T) {
	t.Run("should pass arguments correctly", func(t *testing.T) {
		mockUserID := "mock-user"
		c := mockClient{
			postEphemeral: func(channelID, userID string, options ...slack.MsgOption) (string, error) {
				assert.Len(t, options, 2)
				assert.Equal(t, mockChannelID, channelID)

				return "", nil
			},
		}

		Mute(c, mockChannelID, mockUserID, mockTS)
	})
}

func TestPost(t *testing.T) {
	mockMsg := "mock-message"

	t.Run("should pass arguments correctly", func(t *testing.T) {
		c := mockClient{
			postMessage: func(channelID string, options ...slack.MsgOption) (string, string, error) {
				assert.Len(t, options, 1)
				assert.Equal(t, mockChannelID, channelID)

				return "", "", nil
			},
		}

		Post(c, mockMsg, mockChannelID, mockTS, false)
	})

	t.Run("should apply replyToThread option", func(t *testing.T) {
		c := mockClient{
			postMessage: func(channelID string, options ...slack.MsgOption) (string, string, error) {
				assert.Len(t, options, 2)
				assert.Equal(t, mockChannelID, channelID)

				return "", "", nil
			},
		}

		Post(c, mockMsg, mockChannelID, mockTS, true)
	})
}
