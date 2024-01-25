package test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

type LocalstackContainer struct {
	*localstack.LocalStackContainer
	Port string
}

func NewLocalstackContainer(ctx context.Context) (*LocalstackContainer, error) {
	container, err := localstack.RunContainer(ctx,
		testcontainers.WithImage("localstack/localstack:1.4.0"),
		testcontainers.WithStartupCommand(testcontainers.NewRawCommand([]string{"awslocal ses verify-domain-identity --domain ebook_store.com --endpoint-url http://localhost:4566\""})),
	)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "4566/tcp")
	if err != nil {
		return nil, err
	}

	provider, err := testcontainers.NewDockerProvider()
	if err != nil {
		return nil, err
	}
	defer provider.Close()

	host, err := provider.DaemonHost(ctx)
	if err != nil {
		return nil, err
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, opts ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           fmt.Sprintf("http://%s:%d", host, mappedPort.Int()),
				SigningRegion: region,
			}, nil
		})

	awsConfig, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion("us-east-1"),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "test")),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String("ebook-store"),
	})
	if err != nil {
		return nil, err
	}

	sesClient := ses.NewFromConfig(awsConfig)
	_, err = sesClient.VerifyDomainIdentity(ctx, &ses.VerifyDomainIdentityInput{
		Domain: aws.String("ebook_store.com"),
	})
	if err != nil {
		return nil, err
	}

	return &LocalstackContainer{
		LocalStackContainer: container,
		Port:                mappedPort.Port(),
	}, nil
}
