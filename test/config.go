package test

import (
	"github.com/spf13/viper"
)

func SetEnvironmentVariables() {
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("PORT", "8081")
	viper.SetDefault("ENV", "local")

	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", "5433")
	viper.SetDefault("POSTGRES_USERNAME", "postgres")
	viper.SetDefault("POSTGRES_PASSWORD", "postgres")
	viper.SetDefault("POSTGRES_DATABASE", "postgres")
}
