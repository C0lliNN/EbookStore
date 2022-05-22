package storage_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/c0llinn/ebook-store/internal/config"
	"github.com/c0llinn/ebook-store/internal/storage"
	"github.com/c0llinn/ebook-store/test"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
)

type StorageClientTestSuite struct {
	suite.Suite
	client storage.S3Client
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

	s.client = storage.S3Client{Service: config.NewS3Service(), Bucket: storage.Bucket(viper.GetString("AWS_S3_BUCKET"))}
}

func (s *StorageClientTestSuite) TearDownSuite() {
	ctx := context.TODO()

	s.container.Terminate(ctx)
}

func TestStorageClientTestSuiteRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	suite.Run(t, new(StorageClientTestSuite))
}

func (s *StorageClientTestSuite) TestGeneratePreSignedUrl() {
	key := "some-key"

	url, err := s.client.GeneratePreSignedUrl(context.TODO(), key)

	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), url)
}

func (s *StorageClientTestSuite) TestSaveFile() {
	key := "some-key"
	content := bytes.NewReader([]byte("this is the content of a book"))

	err := s.client.SaveFile(context.TODO(), key, "text/plain", content)

	assert.Nil(s.T(), err)
}

func (s *StorageClientTestSuite) TestRetrieveFile() {
	key := "some-key"
	byts := []byte("this is the content of a book")
	content := bytes.NewReader(byts)

	err := s.client.SaveFile(context.TODO(), key, "text/plain", content)
	assert.Nil(s.T(), err)

	reader, err := s.client.RetrieveFile(context.TODO(), key)
	assert.Nil(s.T(), err)

	actual, err := ioutil.ReadAll(reader)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), byts, actual)
}
