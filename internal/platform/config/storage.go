package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ebookstore/internal/platform/storage"
	"github.com/spf13/viper"
)

func NewBucket() storage.Bucket {
	return storage.Bucket(viper.GetString("AWS_S3_BUCKET"))
}

func NewS3Client(cfg *aws.Config) *s3.Client {
	return s3.NewFromConfig(*cfg, s3.WithEndpointResolver(s3.EndpointResolverFromURL(viper.GetString("AWS_S3_ENDPOINT"))))
}

func NewPresignClient(client *s3.Client) *s3.PresignClient {
	return s3.NewPresignClient(client)
}
