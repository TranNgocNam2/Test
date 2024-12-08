package config

import (
	"github.com/spf13/viper"
	"gitlab.com/innovia69420/kit/enum/file"
	"gitlab.com/innovia69420/kit/enum/message"
	"gitlab.com/innovia69420/kit/logger"
)

type Config struct {
	Host           string `mapstructure:"APP_HOST"`
	Port           int    `mapstructure:"APP_PORT"`
	CorsOrigin     string `mapstructure:"CORS_ORIGIN"`
	CorsProd       string `mapstructure:"CORS_PRODUCTION"`
	DatabaseUrl    string `mapstructure:"DB_DSN"`
	ApiKey         string `mapstructure:"API_KEY"`
	SendGridApiKey string `mapstructure:"SENDGRID_API_KEY"`
	MailName       string `mapstructure:"MAIL_NAME"`
	MAIL_DOMAIN    string `mapstructure:"MAIL_DOMAIN"`
}

func LoadAllAppConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(file.Env)
	viper.AutomaticEnv()
	l := logger.Get(path)

	err = viper.ReadInConfig()
	if err != nil {
		logger.StartUpError(l, message.FailedLoadingEnv)
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		logger.StartUpError(l, message.FailedUnmarshalConfig)
		return
	}
	return
}
