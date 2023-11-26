package fiberx

import (
	"github.com/gofiber/fiber/v2"
	app "github.com/woojae9488/slack-test-bot"
	"github.com/woojae9488/slack-test-bot/fiberx/slack"
)

type Router struct {
	app   fiber.Router
	slack *slack.Controller
}

func NewRouter(
	app *fiber.App,
	config *app.Config,
) *Router {
	return &Router{
		app:   app,
		slack: slack.NewController(config),
	}
}

func (r *Router) Register() {
	slakApi := r.app.Group("/api/slack")
	slakApi.Post("/events", r.slack.AcceptEvent)

	r.app.Use(func(c *fiber.Ctx) error {
		return app.NewError(app.ErrCodeNotFoundPath, c.Path())
	})
}
