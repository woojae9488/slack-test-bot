package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/woojae9488/slack-test-bot/config"
	"github.com/woojae9488/slack-test-bot/handler"
)

var (
	slackHandler *handler.SlackHandler = &handler.SlackHandler{}
	errorHandler *handler.ErrorHandler = &handler.ErrorHandler{}
)

func main() {
	// Create fiber app
	app := fiber.New(fiber.Config{
		AppName: "Slack Test Bot",
		Prefork: config.Server.IsRealPhase(),
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Create a /api/slack endpoint
	slakApi := app.Group("/api/slack")
	// Bind slack api handlers
	slakApi.Post("/events", slackHandler.AcceptEvents)

	// Handle not founds
	app.Use(errorHandler.NotFound)

	// Listen on port 8010
	log.Fatal(app.Listen(config.Server.Port))
}
