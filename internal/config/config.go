package config

import (
	"github.com/spf13/viper"
	"gitlab.com/innovia69420/kit/enum/file"
	"gitlab.com/innovia69420/kit/enum/message"
	"gitlab.com/innovia69420/kit/logger"
)

type Config struct {
	Host        string `mapstructure:"APP_HOST"`
	Port        int    `mapstructure:"APP_PORT"`
	CorsOrigin  string `mapstructure:"CORS_ORIGIN"`
	DatabaseUrl string `mapstructure:"DB_DSN"`
	ApiKey      string `mapstructure:"API_KEY"`
}

func LoadAllAppConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(file.Env)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		l := logger.Get(path)
		logger.StartUpError(l, message.FailedLoadingEnv)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
