package env

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitConfiguration() {
	viper.AutomaticEnv()

	viper.SetConfigName("env")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error %w", err))
	}

	fmt.Printf("Configuration loaded from %s", viper.ConfigFileUsed())
}
