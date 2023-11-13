package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/woojae9488/slack-test-bot/slackextension"
)

var (
	slackSigningSecret = "{SLACK_SIGNING_SECRET}"
	slackClient        = slack.New("{SLACK_ACCES_TOKEN}")
)

const (
	feedbackChannel = "playground"
)

func VerifySlackSecret(reqHeader http.Header, reqBody []byte) *fiber.Error {
	secretsVerifier, err := slack.NewSecretsVerifier(reqHeader, slackSigningSecret)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to make slack secrets verifier")
	}
	if _, err := secretsVerifier.Write(reqBody); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to verify slack secrets")
	}
	if err := secretsVerifier.Ensure(); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Failed to ensure slack secrets")
	}
	return nil
}

func ParseSlackEvent(reqBody []byte) (*slackevents.EventsAPIEvent, *fiber.Error) {
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(reqBody), slackevents.OptionNoVerifyToken())
	if err != nil {
		return &eventsAPIEvent, fiber.NewError(fiber.StatusInternalServerError, "Failed to parse slack event")
	}
	return &eventsAPIEvent, nil
}

func RetrieveSlackEventChallenge(reqBody []byte) (string, *fiber.Error) {
	var challengeResponse *slackevents.ChallengeResponse
	if err := json.Unmarshal(reqBody, &challengeResponse); err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Failed to unmarshal challenge")
	}
	return challengeResponse.Challenge, nil
}

func FeedbackSlackCallbackEvent(innerEvent slackevents.EventsAPIInnerEvent) *fiber.Error {
	switch eventData := innerEvent.Data.(type) {
	case *slackevents.UserProfileChangedEvent:
		if slackevents.EventsAPIType(innerEvent.Type) == slackextension.UserStatusChanged {
			feedbackSlackUserChangedEvent(eventData)
		}
	}
	return nil
}

func feedbackSlackUserChangedEvent(eventData *slackevents.UserProfileChangedEvent) *fiber.Error {
	userName := eventData.User.Name
	userStatus := eventData.User.Profile.StatusText

	var messageText string
	if userStatus == "" {
		messageText = fmt.Sprintf("User `%s` has cleared the status.", userName)
	} else {
		messageText = fmt.Sprintf("User `%s`'s status has changed to `%s`.", userName, userStatus)
	}

	channelId, timestamp, err := slackClient.PostMessage(feedbackChannel, slack.MsgOptionText(messageText, false))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to post feedback message")
	}
	log.Infof("[SLACK][EVENT] Post feedback message to %s at %s", channelId, timestamp)

	return nil
}
