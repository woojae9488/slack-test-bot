package main

import (
	app "github.com/woojae9488/slack-test-bot"
	"github.com/woojae9488/slack-test-bot/fiberx"
	"github.com/woojae9488/slack-test-bot/viperx"
)

func main() {
	c := viperx.NewConfig()
	s := fiberx.NewServer(c)
	app.StartServer(s)
}
