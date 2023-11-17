//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/woojae9488/slack-test-bot/config"
	"github.com/woojae9488/slack-test-bot/handler"
	"github.com/woojae9488/slack-test-bot/slack"
)

func initializeServer() *Server {
	wire.Build(
		config.NewServerConfig,
		config.NewSlackConfig,
		slack.NewSlackClient,
		slack.NewEventService,
		handler.NewSlackHandler,
		handler.NewErrorHandler,
		NewFiberApp,
		NewServer,
	)
	return &Server{}
}
