package test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
	viper.SetDefault("AWS_S3_BUCKET", "ebook-store")
	viper.SetDefault("AWS_S3_ENDPOINT", "http://s3.localhost.localstack.cloud:5566")

	viper.SetDefault("STRIPE_API_KEY", "sk_test_51HAKIGHKmAtjDhlfifsr2lIoY8nQZXkQTE2RvqFfa4ASe6Rlk4YRfVxp44Rr9eeSrPivk55dloy9KFv5Zal3sWQz009q9hiu1u")
}

func NewS3Service() *s3.S3 {
	var endpoint *string
	if env := viper.GetString("AWS_S3_ENDPOINT"); env != "" {
		endpoint = aws.String(env)
	}

	region := viper.GetString("AWS_REGION")
	currentSession, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: endpoint,
	})

	if err != nil {
		panic(err)
	}

	return s3.New(currentSession)
}
