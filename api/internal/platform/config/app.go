package config

import (
	"Backend/kit/enum"
	"Backend/kit/log"
	"github.com/spf13/viper"
)

type App struct {
	Host       string `mapstructure:"APP_HOST"`
	Port       int    `mapstructure:"APP_PORT"`
	CorsOrigin string `mapstructure:"CORS_ORIGIN"`
}

func LoadAllAppConfig(path string) (app *App, err error) {
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
