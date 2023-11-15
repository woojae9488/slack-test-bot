package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/slack-go/slack/slackevents"
	"github.com/woojae9488/slack-test-bot/slack"
)

var (
	eventService *slack.EventService = &slack.EventService{}
)

type SlackHandler struct{}

func (handler *SlackHandler) AcceptEvents(c *fiber.Ctx) error {
	header := http.Header(c.GetReqHeaders())
	body := c.Body()

	if err := eventService.VerifySecret(header, body); err != nil {
		return err
	}

	eventsAPIEvent, err := eventService.ParseEvent(body)
	if err != nil {
		return err
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		challenge, err := eventService.RetrieveEventChallenge(body)
		if err != nil {
			return err
		}
		return c.Type(fiber.MIMETextPlain).Send([]byte(challenge))
	} else if eventsAPIEvent.Type == slackevents.CallbackEvent {
		if err := eventService.FeedbackCallbackEvent(eventsAPIEvent.InnerEvent); err != nil {
			return err
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}
