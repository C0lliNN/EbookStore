package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ebookstore/internal/log"
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

func (c *Storage) GenerateGetPreSignedUrl(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	presignResult, err := c.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(string(c.Bucket)),
		Key:    aws.String(key),
	})

	if err != nil {
		return "", fmt.Errorf("(GenerateGetPreSignedUrl) failed generating presignedUrl for key: %s: %w", key, err)
	}

	return presignResult.URL, nil
}

func (c *Storage) GeneratePutPreSignedUrl(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	log.FromContext(ctx).Info("generating PUT presigned url for key: %s", key)

	presignResult, err := c.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(string(c.Bucket)),
		Key:    aws.String(key),
	})

	if err != nil {
		return "", fmt.Errorf("(GeneratePutPreSignedUrl) failed generating presignedUrl for key: %s: %w", key, err)
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
		return fmt.Errorf("(SaveFile) failed saving file for key %s: %w", key, err)
	}

	return nil
}

func (c *Storage) RetrieveFile(ctx context.Context, key string) (io.ReadCloser, error) {
	output, err := c.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(string(c.Bucket)),
	})
	if err != nil {
		return nil, fmt.Errorf("(RetrieveFile) failed retrieving file for key %s: %w", key, err)
	}

	return output.Body, nil
}
