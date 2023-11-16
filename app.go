package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/woojae9488/slack-test-bot/config"
	"github.com/woojae9488/slack-test-bot/handler"
)

type App struct {
	config config.Config
	slackH handler.SlackHandler
	errorH handler.ErrorHandler
}

func NewApp(slackH handler.SlackHandler, errorH handler.ErrorHandler) App {
	return App{
		slackH: slackH,
		errorH: errorH,
	}
}

func main() {
	a := initializeApp()

	// Create fiber app
	app := fiber.New(fiber.Config{
		AppName: "Slack Test Bot",
		Prefork: a.config.Server.IsRealPhase(),
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Create a /api/slack endpoint
	slakApi := app.Group("/api/slack")
	// Bind slack api handlers
	slakApi.Post("/events", a.slackH.AcceptEvents)

	// Handle not founds
	app.Use(a.errorH.NotFound)

	// Listen on port 8010
	log.Fatal(app.Listen(a.config.Server.Port))
}
