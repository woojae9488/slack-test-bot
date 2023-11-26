package slack

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	app "github.com/woojae9488/slack-test-bot"
)

var (
	UserStatusChanged = slackevents.EventsAPIType("user_status_changed")
)

func init() {
	slackevents.EventsAPIInnerEventMapping[UserStatusChanged] = slackevents.UserProfileChangedEvent{}
}

func NewClient(config *app.Config) *slack.Client {
	return slack.New(config.SlackToken)
}
