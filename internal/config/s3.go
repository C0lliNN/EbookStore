package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/c0llinn/ebook-store/internal/log"
	"github.com/c0llinn/ebook-store/internal/storage"
	"github.com/spf13/viper"
)

func NewBucket() storage.Bucket {
	return storage.Bucket(viper.GetString("AWS_S3_BUCKET"))
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
		log.Default().Fatalf("could not create an aws session: %v", err)
	}

	return s3.New(currentSession)
}
