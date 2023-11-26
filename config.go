package slacktestbot

import "strconv"

type ServerPhase string

const (
	ServerLocalPhase ServerPhase = "local"
	ServerRealPhase  ServerPhase = "real"
)

type Config struct {
	ServerPhase ServerPhase
	ServerPort  int

	SlackToken           string
	SlackSigningSecret   string
	SlackFeedbackChannel string
}

func (c *Config) IsRealPhase() bool {
	return c.ServerPhase == ServerRealPhase
}

func (c *Config) GetAddress() string {
	return ":" + strconv.Itoa(c.ServerPort)
}

var ConfigDefault = Config{
	ServerPhase: ServerLocalPhase,
	ServerPort:  8010,

	SlackToken:           "",
	SlackSigningSecret:   "",
	SlackFeedbackChannel: "playground",
}
