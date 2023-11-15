package slack

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/woojae9488/slack-test-bot/config"
)

var (
	client          = slack.New(config.Slack.Token)
	signingSecret   = config.Slack.SigningSecret
	feedbackChannel = config.Slack.FeedbackChannel
)

type EventService struct{}

func (service *EventService) VerifySecret(reqHeader http.Header, reqBody []byte) error {
	secretsVerifier, err := slack.NewSecretsVerifier(reqHeader, signingSecret)
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

func (service *EventService) ParseEvent(reqBody []byte) (*slackevents.EventsAPIEvent, error) {
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(reqBody), slackevents.OptionNoVerifyToken())
	if err != nil {
		return &eventsAPIEvent, fiber.NewError(fiber.StatusInternalServerError, "Failed to parse slack event")
	}
	return &eventsAPIEvent, nil
}

func (service *EventService) RetrieveEventChallenge(reqBody []byte) (string, error) {
	var challengeResponse *slackevents.ChallengeResponse
	if err := json.Unmarshal(reqBody, &challengeResponse); err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Failed to unmarshal slack challenge")
	}
	return challengeResponse.Challenge, nil
}

func (service *EventService) FeedbackCallbackEvent(innerEvent slackevents.EventsAPIInnerEvent) error {
	eventsApiType := slackevents.EventsAPIType(innerEvent.Type)
	switch eventData := innerEvent.Data.(type) {
	case *slackevents.UserProfileChangedEvent:
		if eventsApiType == UserStatusChanged {
			service.feedbackUserChangedEvent(eventData)
		}
	}
	return nil
}

func (service *EventService) feedbackUserChangedEvent(eventData *slackevents.UserProfileChangedEvent) error {
	userName := eventData.User.Name
	userStatus := eventData.User.Profile.StatusText

	var message string
	if userStatus == "" {
		message = fmt.Sprintf("User `%s` has cleared the status.", userName)
	} else {
		message = fmt.Sprintf("User `%s`'s status has changed to `%s`.", userName, userStatus)
	}

	if err := service.postFeedbackMessage(message); err != nil {
		return err
	}
	return nil
}

func (service *EventService) postFeedbackMessage(message string) error {
	channelId, timestamp, err := client.PostMessage(feedbackChannel, slack.MsgOptionText(message, false))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to post slack feedback message")
	}

	log.Infof("[SLACK][EVENT] Post feedback message to %s at %s", channelId, timestamp)
	return nil
}
