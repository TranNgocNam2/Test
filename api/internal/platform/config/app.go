package config

import (
	"Backend/kit/enum"
	"Backend/kit/log"
	"github.com/spf13/viper"
)

type AppCfg struct {
	Host        string `mapstructure:"APP_HOST"`
	Port        int    `mapstructure:"APP_PORT"`
	CorsOrigin  string `mapstructure:"CORS_ORIGIN"`
	DatabaseUrl string `mapstructure:"DATABASE_URL"`
}

var App *AppCfg

func LoadAllAppConfig(path string) (app *AppCfg, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(enum.EnvironmentFile)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		l := log.Get(path)
		log.StartUpError(l, enum.ErrorLoadEnvFile)
		return
	}

	err = viper.Unmarshal(&app)
	return
}
