package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/spf13/viper"
)

func NewSNSService() *ses.SES {
	var endpoint *string
	if env := viper.GetString("AWS_SES_ENDPOINT"); env != "" {
		endpoint = aws.String(env)
	}

	region := viper.GetString("AWS_REGION")

	currentSession, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: endpoint,
	})

	if err != nil {
		log.Logger.Fatalw("Error creating aws session", "error", err)
	}

	return ses.New(currentSession)
}
