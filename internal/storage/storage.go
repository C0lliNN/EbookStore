package storage

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/c0llinn/ebook-store/internal/log"
	"io"
	"time"
)

type Config struct {
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
	Bucket        Bucket
}

type Bucket string

type Storage struct {
	Config
}

func NewStorage(c Config) *Storage {
	return &Storage{Config: c}
}

func (c *Storage) GeneratePreSignedUrl(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	presignResult, err := c.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(string(c.Bucket)),
		Key:    aws.String(key),
	})

	if err != nil {
		log.FromContext(ctx).Errorf("error generating get presignUrl for key %s: %v", key, err)
		return "", err
	}

	return presignResult.URL, nil
}

func (c *Storage) SaveFile(ctx context.Context, key string, contentType string, content io.ReadSeeker) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	_, err := c.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Key:         aws.String(key),
		Bucket:      aws.String(string(c.Bucket)),
		ContentType: aws.String(contentType),
		Body:        content,
	})
	if err != nil {
		log.FromContext(ctx).Errorf("error saving file for key %s: %v", key, err)
		return err
	}

	return nil
}

func (c *Storage) RetrieveFile(ctx context.Context, key string) (io.ReadCloser, error) {
	output, err := c.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(string(c.Bucket)),
	})
	if err != nil {
		log.FromContext(ctx).Errorf("error retriving file for key %s: %v", key, err)
		return nil, err
	}

	return output.Body, nil
}
