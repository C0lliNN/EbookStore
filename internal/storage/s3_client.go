package storage

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/c0llinn/ebook-store/internal/log"
	"io"
	"time"
)

type Bucket string

type S3Client struct {
	service *s3.S3
	bucket  Bucket
}

func NewS3Client(service *s3.S3, bucket Bucket) *S3Client {
	return &S3Client{service: service, bucket: bucket}
}

func (c *S3Client) GeneratePreSignedUrl(ctx context.Context, key string) (string, error) {
	request, _ := c.service.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(string(c.bucket)),
		Key:    aws.String(key),
	})

	url, err := request.Presign(time.Hour)
	if err != nil {
		log.Default().Errorf("Error generating get presignUrl for key %s: %v", key, err)
	}

	return url, err
}

func (c *S3Client) SaveFile(ctx context.Context, key string, contentType string, content io.ReadSeeker) error {
	_, err := c.service.PutObject(&s3.PutObjectInput{
		Key:         aws.String(key),
		Bucket:      aws.String(string(c.bucket)),
		ContentType: aws.String(contentType),
		Body:        content,
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *S3Client) RetrieveFile(ctx context.Context, key string) (io.ReadCloser, error) {
	output, err := c.service.GetObject(&s3.GetObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(string(c.bucket)),
	})

	if err != nil {
		return nil, err
	}

	return output.Body, nil
}
