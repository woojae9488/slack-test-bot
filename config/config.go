package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/woojae9488/slack-test-bot/util"
)

type ServerPhase string

const (
	ServerLocalPhase ServerPhase = "local"
	ServerRealPhase  ServerPhase = "real"

	ServerDefaultPhase ServerPhase = ServerLocalPhase
	ServerDefaultPort  string      = ":8010"

	AppConfigDir      string = "config"
	AppConfigPrefix   string = "app"
	SlackConfigPrefix string = AppConfigPrefix + ".slack"
)

type ServerConfig struct {
	Port  string
	Phase ServerPhase
}

type SlackConfig struct {
	Token           string `mapstructure:"token"`
	SigningSecret   string `mapstructure:"signing-secret"`
	FeedbackChannel string `mapstructure:"feedback-channel"`
}

func (s *ServerConfig) IsRealPhase() bool {
	return s.Phase == ServerRealPhase
}

func NewServerConfig() *ServerConfig {
	pflag.String("port", ServerDefaultPort, "Port to listen on")
	pflag.String("phase", string(ServerDefaultPhase), "Enable prefork on real phase")
	pflag.Parse()
	util.Validate(viper.BindPFlags(pflag.CommandLine))

	s := ServerConfig{}
	util.Validate(viper.Unmarshal(&s))
	return &s
}

func NewSlackConfig(server *ServerConfig) *SlackConfig {
	viper.AddConfigPath(AppConfigDir)
	viper.SetConfigName(string(server.Phase))
	util.Validate(viper.ReadInConfig())

	s := SlackConfig{}
	util.Validate(viper.Sub(SlackConfigPrefix).Unmarshal(&s))
	return &s
}
