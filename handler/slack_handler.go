package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/slack-go/slack/slackevents"
	"github.com/woojae9488/slack-test-bot/service"
)

func SlackEvents(c *fiber.Ctx) error {
	header := http.Header(c.GetReqHeaders())
	body := c.Body()

	if err := service.VerifySlackSecret(header, body); err != nil {
		return err
	}

	eventsAPIEvent, err := service.ParseSlackEvent(body)
	if err != nil {
		return err
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		challenge, err := service.RetrieveSlackEventChallenge(body)
		if err != nil {
			return err
		}
		return c.Type(fiber.MIMETextPlain).Send([]byte(challenge))
	} else if eventsAPIEvent.Type == slackevents.CallbackEvent {
		if err := service.FeedbackSlackCallbackEvent(eventsAPIEvent.InnerEvent); err != nil {
			return err
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}
