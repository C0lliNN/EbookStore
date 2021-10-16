//go:build integration
// +build integration

package storage

import (
	"bytes"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/c0llinn/ebook-store/test"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"testing"
)

type StorageClientTestSuite struct {
	suite.Suite
	client S3Client
}

func (s *StorageClientTestSuite) SetupTest() {
	test.SetEnvironmentVariables()
	log.InitLogger()

	s.client = S3Client{service: test.NewS3Service(), bucket: Bucket(viper.GetString("AWS_S3_BUCKET"))}
}

func TestStorageClientTestSuiteRun(t *testing.T) {
	suite.Run(t, new(StorageClientTestSuite))
}

func (s *StorageClientTestSuite) TestGeneratePreSignedUrl() {
	key := "some-key"

	url, err := s.client.GeneratePreSignedUrl(key)

	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), url)
}

func (s *StorageClientTestSuite) TestSaveFile() {
	key := "some-key"
	content := bytes.NewReader([]byte("this is the content of a book"))

	err := s.client.SaveFile(key, "text/plain", content)

	assert.Nil(s.T(), err)
}

func (s *StorageClientTestSuite) TestRetrieveFile() {
	key := "some-key"
	byts := []byte("this is the content of a book")
	content := bytes.NewReader(byts)

	err := s.client.SaveFile(key, "text/plain", content)
	assert.Nil(s.T(), err)

	reader, err := s.client.RetrieveFile(key)
	assert.Nil(s.T(), err)

	actual, err := ioutil.ReadAll(reader)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), byts, actual)
}
