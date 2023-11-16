package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

type Config struct {
	Server ServerConfig
	Slack  SlackConfig
}

type ServerConfig struct {
	Port  string
	Phase ServerPhase
}

type SlackConfig struct {
	Token           string `mapstructure:"token"`
	SigningSecret   string `mapstructure:"signing-secret"`
	FeedbackChannel string `mapstructure:"feedback-channel"`
}

func NewConfig(server ServerConfig, slack SlackConfig) Config {
	return Config{
		Server: server,
		Slack:  slack,
	}
}

func (s *ServerConfig) IsRealPhase() bool {
	return s.Phase == ServerRealPhase
}

func NewServerConfig() ServerConfig {
	pflag.String("port", ServerDefaultPort, "Port to listen on")
	pflag.String("phase", string(ServerDefaultPhase), "Enable prefork on real phase")
	pflag.Parse()
	validate(viper.BindPFlags(pflag.CommandLine))

	s := ServerConfig{}
	validate(viper.Unmarshal(&s))
	return s
}

func NewSlackConfig(server ServerConfig) SlackConfig {
	viper.AddConfigPath(AppConfigDir)
	viper.SetConfigName(string(server.Phase))
	validate(viper.ReadInConfig())

	s := SlackConfig{}
	validate(viper.Sub(SlackConfigPrefix).Unmarshal(&s))
	return s
}

func validate(err error) {
	if err != nil {
		panic(err)
	}
}
