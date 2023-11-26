package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	app "github.com/woojae9488/slack-test-bot"
)

type EventService struct {
	client          *slack.Client
	signingSecret   string
	feedbackChannel string
}

func NewEventService(client *slack.Client, config *app.Config) *EventService {
	return &EventService{
		client:          client,
		signingSecret:   config.SlackSigningSecret,
		feedbackChannel: config.SlackFeedbackChannel,
	}
}

func (s *EventService) VerifySecret(reqHeader http.Header, reqBody []byte) error {
	secretsVerifier, err := slack.NewSecretsVerifier(reqHeader, s.signingSecret)
	if err != nil {
		return errors.Join(fiber.NewError(fiber.StatusBadRequest, "Http headers or signging secret is invalid"), err)
	}
	if _, err := secretsVerifier.Write(reqBody); err != nil {
		return errors.Join(fiber.NewError(fiber.StatusInternalServerError, "Failed to write verifier hmac"), err)
	}
	if err := secretsVerifier.Ensure(); err != nil {
		return errors.Join(fiber.NewError(fiber.StatusUnauthorized, "Failed to ensure secret hmac"), err)
	}
	return nil
}

func (s *EventService) ParseEvent(reqBody []byte) (*slackevents.EventsAPIEvent, error) {
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(reqBody), slackevents.OptionNoVerifyToken())
	if err != nil {
		return &eventsAPIEvent,
			errors.Join(fiber.NewError(fiber.StatusInternalServerError, "Failed to parse event"), err)
	}
	return &eventsAPIEvent, nil
}

func (s *EventService) RetrieveEventChallenge(reqBody []byte) (string, error) {
	var challengeResponse *slackevents.ChallengeResponse
	if err := json.Unmarshal(reqBody, &challengeResponse); err != nil {
		return "",
			errors.Join(fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve event challenge"), err)
	}
	return challengeResponse.Challenge, nil
}

func (s *EventService) FeedbackCallbackEvent(innerEvent slackevents.EventsAPIInnerEvent) error {
	eventsApiType := slackevents.EventsAPIType(innerEvent.Type)
	switch eventData := innerEvent.Data.(type) {
	case *slackevents.UserProfileChangedEvent:
		if eventsApiType == UserStatusChanged {
			return s.feedbackUserChangedEvent(eventData)
		}
	}
	return nil
}

func (s *EventService) feedbackUserChangedEvent(eventData *slackevents.UserProfileChangedEvent) error {
	userName := eventData.User.Name
	userStatus := eventData.User.Profile.StatusText

	var message string
	if userStatus == "" {
		message = fmt.Sprintf("User `%s` has cleared the status.", userName)
	} else {
		message = fmt.Sprintf("User `%s`'s status has changed to `%s`.", userName, userStatus)
	}

	return s.postFeedbackMessage(message)
}

func (s *EventService) postFeedbackMessage(message string) error {
	channelId, timestamp, err := s.client.PostMessage(s.feedbackChannel, slack.MsgOptionText(message, false))
	if err != nil {
		return errors.Join(fiber.NewError(fiber.StatusInternalServerError, "Failed to post feedback message"), err)
	}

	log.Infof("[SLACK][EVENT] Post feedback message to %s at %s", channelId, timestamp)
	return nil
}
