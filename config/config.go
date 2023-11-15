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

type ServerConfig struct {
	Port  string
	Phase ServerPhase
}

type SlackConfig struct {
	Token           string `mapstructure:"token"`
	SigningSecret   string `mapstructure:"signing-secret"`
	FeedbackChannel string `mapstructure:"feedback-channel"`
}

var (
	Server ServerConfig
	Slack  SlackConfig
)

func (server *ServerConfig) IsRealPhase() bool {
	return server.Phase == ServerRealPhase
}

func init() {
	// Parse command-line flags
	pflag.String("port", ServerDefaultPort, "Port to listen on")
	pflag.String("phase", string(ServerDefaultPhase), "Enable prefork on real phase")
	pflag.Parse()
	validate(viper.BindPFlags(pflag.CommandLine))
	validate(viper.Unmarshal(&Server))

	// Parse yaml properties
	viper.AddConfigPath(AppConfigDir)
	viper.SetConfigName(string(Server.Phase))
	validate(viper.ReadInConfig())
	validate(viper.Sub(SlackConfigPrefix).Unmarshal(&Slack))
}

func validate(err error) {
	if err != nil {
		panic(err)
	}
}
