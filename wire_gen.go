// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/woojae9488/slack-test-bot/config"
	"github.com/woojae9488/slack-test-bot/handler"
	"github.com/woojae9488/slack-test-bot/slack"
)

// Injectors from wire.go:

func initializeServer() *Server {
	serverConfig := config.NewServerConfig()
	app := NewFiberApp(serverConfig)
	slackConfig := config.NewSlackConfig(serverConfig)
	client := slack.NewSlackClient(slackConfig)
	eventService := slack.NewEventService(client, slackConfig)
	slackHandler := handler.NewSlackHandler(eventService)
	errorHandler := handler.NewErrorHandler()
	server := NewServer(app, serverConfig, slackHandler, errorHandler)
	return server
}
