package slack

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/slack-go/slack/slackevents"
	app "github.com/woojae9488/slack-test-bot"
)

type Controller struct {
	event *EventService
}

func NewController(config *app.Config) *Controller {
	client := NewClient(config)
	event := NewEventService(client, config)
	return &Controller{event: event}
}

func (sc *Controller) AcceptEvent(c *fiber.Ctx) error {
	header := http.Header(c.GetReqHeaders())
	body := c.Body()

	if err := sc.event.VerifySecret(header, body); err != nil {
		return err
	}

	eventsAPIEvent, err := sc.event.ParseEvent(body)
	if err != nil {
		return err
	}

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		challenge, err := sc.event.RetrieveEventChallenge(body)
		if err != nil {
			return err
		}
		return c.SendString(challenge)
	case slackevents.CallbackEvent:
		if err := sc.event.FeedbackCallbackEvent(eventsAPIEvent.InnerEvent); err != nil {
			return err
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}
