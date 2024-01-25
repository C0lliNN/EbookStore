package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ebookstore/internal/platform/config"
	"github.com/ebookstore/internal/platform/storage"
	"github.com/ebookstore/test"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StorageClientTestSuite struct {
	suite.Suite
	storage   *storage.Storage
	container *test.LocalstackContainer
}

func (s *StorageClientTestSuite) SetupSuite() {
	config.LoadConfiguration()

	ctx := context.TODO()

	var err error
	s.container, err = test.NewLocalstackContainer(ctx)
	s.Require().NoError(err)

	viper.Set("AWS_S3_ENDPOINT", fmt.Sprintf("http://s3.localhost.localstack.cloud:%v", s.container.Port))

	s.storage = storage.NewStorage(storage.Config{
		S3Client:      config.NewS3Client(config.NewAWSConfig()),
		PresignClient: config.NewPresignClient(config.NewS3Client(config.NewAWSConfig())),
		Bucket:        storage.Bucket(viper.GetString("AWS_S3_BUCKET")),
	})
}

func (s *StorageClientTestSuite) TearDownSuite() {
	ctx := context.TODO()

	_ = s.container.Terminate(ctx)
}

func TestStorageClientTestSuiteRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	suite.Run(t, new(StorageClientTestSuite))
}

func (s *StorageClientTestSuite) TestGenerateGetPreSignedUrl() {
	key := "some-key"

	url, err := s.storage.GenerateGetPreSignedUrl(context.TODO(), key)

	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), url)
}

func (s *StorageClientTestSuite) TestGeneratePutPreSignedUrl() {
	key := "some-key"

	url, err := s.storage.GeneratePutPreSignedUrl(context.TODO(), key)

	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), url)
}
