package fiberx

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	app "github.com/woojae9488/slack-test-bot"
)

type Server struct {
	app    *fiber.App
	config *app.Config
}

func NewServer(
	config *app.Config,
) *Server {
	return &Server{
		app:    newFiberApp(config),
		config: config,
	}
}

func newFiberApp(config *app.Config) *fiber.App {
	return fiber.New(fiber.Config{
		AppName:      "Slack Test Bot",
		Prefork:      config.IsRealPhase(),
		ErrorHandler: ErrorHandler,
	})
}

func (s *Server) SetupMiddlewares() {
	s.app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	s.app.Use(logger.New())
	s.app.Get("/metrics", monitor.New())
}

func (s *Server) SetupRoutes() {
	router := NewRouter(s.app, s.config)
	router.Register()
}

func (s *Server) Listen() {
	log.Fatal(s.app.Listen(s.config.GetAddress()))
}
