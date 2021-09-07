// +build integration

package repository

import (
	"github.com/c0llinn/ebook-store/config/db"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/c0llinn/ebook-store/test"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	repo UserRepository
}

func (s *UserRepositoryTestSuite) SetupTest() {
	test.SetEnvironmentVariables()
	log.InitLogger()

	conn := db.NewConnection()
	s.repo = UserRepository{conn}
}

func TestUserRepositoryTest(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) TestUserRepository_SaveSuccessfully() {
	user := factory.NewUser()

	err := s.repo.Save(&user)

	assert.Nil(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestUserRepository_WithError() {
	user := factory.NewUser()
	user.ID += user.ID

	err := s.repo.Save(&user)

	assert.NotNil(s.T(), err)
}
