package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func InitConfiguration() {
	if os.Getenv("ENV") == "production" {
		viper.AutomaticEnv()

		fmt.Println("Configuration loaded from environment variables")
	} else {
		viper.AddConfigPath("..")
		viper.SetConfigName("env")

		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error %w", err))
		}

		fmt.Printf("Configuration loaded from %s", viper.ConfigFileUsed())
	}
}
