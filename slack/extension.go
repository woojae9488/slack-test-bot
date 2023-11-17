package slack

import (
	"github.com/slack-go/slack/slackevents"
)

var (
	UserStatusChanged = slackevents.EventsAPIType("user_status_changed")
)

func init() {
	slackevents.EventsAPIInnerEventMapping[UserStatusChanged] = slackevents.UserProfileChangedEvent{}
}
