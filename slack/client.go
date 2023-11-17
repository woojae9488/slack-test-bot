package slack

import (
	"github.com/slack-go/slack"
	"github.com/woojae9488/slack-test-bot/config"
)

func NewSlackClient(config *config.SlackConfig) *slack.Client {
	return slack.New(config.Token)
}
