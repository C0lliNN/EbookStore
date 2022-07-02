package storage_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ebookstore/internal/config"
	"github.com/ebookstore/internal/storage"
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
	if err != nil {
		s.T().Fatal(err)
	}

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

func (s *StorageClientTestSuite) TestGeneratePreSignedUrl() {
	key := "some-key"

	url, err := s.storage.GeneratePreSignedUrl(context.TODO(), key)

	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), url)
}

func (s *StorageClientTestSuite) TestSaveFile() {
	key := "some-key"
	content := bytes.NewReader([]byte("this is the content of a book"))

	err := s.storage.SaveFile(context.TODO(), key, "text/plain", content)

	assert.Nil(s.T(), err)
}

func (s *StorageClientTestSuite) TestRetrieveFile() {
	key := "some-key"
	byts := []byte("this is the content of a book")
	content := bytes.NewReader(byts)

	err := s.storage.SaveFile(context.TODO(), key, "text/plain", content)
	assert.Nil(s.T(), err)

	reader, err := s.storage.RetrieveFile(context.TODO(), key)
	assert.Nil(s.T(), err)

	actual, err := ioutil.ReadAll(reader)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), byts, actual)
}
