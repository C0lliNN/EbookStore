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

	viper.SetDefault("AWS_REGION", "us-east-2")
	viper.SetDefault("AWS_SES_ENDPOINT", "http://localhost:5566")
	viper.SetDefault("AWS_SES_SOURCE_EMAIL", "no-reply@ebook_store.com")
}
