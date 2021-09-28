// +build integration

package storage

import (
	"bytes"
	"fmt"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/c0llinn/ebook-store/test"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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

	err := s.client.SaveFile(key, content)

	fmt.Println(err)
	assert.Nil(s.T(), err)
}
