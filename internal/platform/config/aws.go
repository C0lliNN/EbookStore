package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/ebookstore/internal/log"
	"github.com/spf13/viper"
)

func NewAWSConfig() *aws.Config {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(viper.GetString("AWS_REGION")))
	if err != nil {
		log.Fatalf(ctx, "failed to load AWS configuration: %v", err)
	}

	return &cfg
}
