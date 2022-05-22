package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/spf13/viper"
)

func NewSESClient(cfg *aws.Config) *ses.Client {
	return ses.NewFromConfig(*cfg, ses.WithEndpointResolver(ses.EndpointResolverFromURL(viper.GetString("AWS_SES_ENDPOINT"))))
}
