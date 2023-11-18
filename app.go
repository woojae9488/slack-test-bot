package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/woojae9488/slack-test-bot/config"
	"github.com/woojae9488/slack-test-bot/handler"
)

type Server struct {
	app          *fiber.App
	config       *config.ServerConfig
	slackHandler *handler.SlackHandler
	errorHandler *handler.ErrorHandler
}

func NewServer(
	app *fiber.App,
	config *config.ServerConfig,
	slackHandler *handler.SlackHandler,
	errorHandler *handler.ErrorHandler,
) *Server {
	return &Server{
		app:          app,
		config:       config,
		slackHandler: slackHandler,
		errorHandler: errorHandler,
	}
}

func NewFiberApp(config *config.ServerConfig) *fiber.App {
	return fiber.New(fiber.Config{
		AppName: "Slack Test Bot",
		Prefork: config.IsRealPhase(),
	})
}

func (s *Server) setupMiddlewares() {
	// Middleware
	s.app.Use(recover.New())
	s.app.Use(logger.New())
}

func (s *Server) setupHandlers() {
	// Create a /api/slack endpoint
	slakApi := s.app.Group("/api/slack")
	// Bind slack api handlers
	slakApi.Post("/events", s.slackHandler.AcceptEvent)

	// Handle not founds
	s.app.Use(s.errorHandler.NotFound)
}

func (s *Server) startListen() {
	// Listen on port
	log.Fatal(s.app.Listen(s.config.Port))
}

func main() {
	s := initializeServer()
	s.setupMiddlewares()
	s.setupHandlers()
	s.startListen()
}
