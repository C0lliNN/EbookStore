package env

import (
	"fmt"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/spf13/viper"
	"os"
)

func InitConfiguration() {
	if os.Getenv("ENV") == "production" {
		viper.AutomaticEnv()

		log.Logger.Debug("Configuration loaded from environment variables")
	} else {
		viper.AddConfigPath("..")
		viper.SetConfigName("env")

		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error %w", err))
		}

		log.Logger.Debug("Configuration loaded from %s", viper.ConfigFileUsed())
	}
}
