package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

func LoadConfiguration() {
	switch os.Getenv("ENV") {
	case "production":
		{
			viper.AutomaticEnv()

			log.Println("configuration loaded from environment variables")
			return
		}
	case "test":
		{
			viper.AddConfigPath("../../..")
			viper.SetConfigName("env-test")
		}
	default: // local is the default
		viper.AddConfigPath("..")
		viper.SetConfigName("env-local")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	log.Println("configuration loaded from ", viper.ConfigFileUsed())
}
