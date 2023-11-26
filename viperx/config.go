package viperx

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	app "github.com/woojae9488/slack-test-bot"
)

var (
	configDir               = "."
	slakConfigPrefix string = "app.slack"
)

type serverConfig struct {
	Phase app.ServerPhase `mapstructure:"phase"`
	Port  int             `mapstructure:"port"`
}

type slackConfig struct {
	Token           string `mapstructure:"token"`
	SigningSecret   string `mapstructure:"signing-secret"`
	FeedbackChannel string `mapstructure:"feedback-channel"`
}

func NewConfig() *app.Config {
	server := newServerConfig()
	slack := newSlackConfig(server.Phase)

	return &app.Config{
		ServerPhase:          server.Phase,
		ServerPort:           server.Port,
		SlackToken:           slack.Token,
		SlackSigningSecret:   slack.SigningSecret,
		SlackFeedbackChannel: slack.FeedbackChannel,
	}
}

func newServerConfig() serverConfig {
	pflag.String("phase", string(app.ConfigDefault.ServerPhase), "Enable prefork on real phase")
	pflag.Int("port", app.ConfigDefault.ServerPort, "Port to listen on")
	pflag.Parse()
	app.PanicOnError(viper.BindPFlags(pflag.CommandLine))

	s := serverConfig{}
	app.PanicOnError(viper.Unmarshal(&s))
	return s
}

func newSlackConfig(phase app.ServerPhase) slackConfig {
	viper.AddConfigPath(configDir)
	viper.SetConfigName(string(phase))
	app.PanicOnError(viper.ReadInConfig())

	s := slackConfig{
		Token:           app.ConfigDefault.SlackToken,
		SigningSecret:   app.ConfigDefault.SlackSigningSecret,
		FeedbackChannel: app.ConfigDefault.SlackFeedbackChannel,
	}
	app.PanicOnError(viper.Sub(slakConfigPrefix).Unmarshal(&s))
	return s
}
