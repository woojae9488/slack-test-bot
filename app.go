package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/woojae9488/slack-welcome-bot/handler"
)

var (
	port    = flag.String("port", ":8010", "Port to listen on")
	profile = flag.String("profile", "local", "Enable prefork on real profile")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Create fiber app
	app := fiber.New(fiber.Config{
		AppName: "Slack Test Bot",
		Prefork: *profile == "real",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Create a /api/slack endpoint
	slakApi := app.Group("/api/slack")
	// Bind slack api handlers
	slakApi.Post("/events", handler.SlackEvents)

	// Handle not founds
	app.Use(handler.NotFound)

	// Listen on port 8010
	log.Fatal(app.Listen(*port))
}
