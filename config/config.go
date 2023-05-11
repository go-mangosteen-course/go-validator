package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadAppConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config/mangosteen/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
}
