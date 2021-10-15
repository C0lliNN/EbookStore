//go:build integration
// +build integration

package repository

import (
	"github.com/c0llinn/ebook-store/config/db"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/c0llinn/ebook-store/test"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	db.LoadMigrations("file:../../../migration")

	conn := db.NewConnection()
	s.repo = UserRepository{conn}
}

func (s *UserRepositoryTestSuite) TearDownTest() {
	s.repo.db.Delete(&model.User{}, "1 = 1")
}

func TestUserRepositoryTest(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) TestUserRepository_SaveSuccessfully() {
	user := factory.NewUser()

	err := s.repo.Save(&user)

	assert.Nil(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestUserRepository_SaveWithDuplicateEmail() {
	user := factory.NewUser()

	err := s.repo.Save(&user)
	assert.Nil(s.T(), err)

	err = s.repo.Save(&user)
	assert.IsType(s.T(), &common.ErrDuplicateKey{}, err)
}

func (s *UserRepositoryTestSuite) TestUserRepository_FindByEmailSuccessfully() {
	expected := factory.NewUser()

	err := s.repo.Save(&expected)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByEmail(expected.Email)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)
}

func (s *UserRepositoryTestSuite) TestUserRepository_FindByEmailNotFound() {
	_, err := s.repo.FindByEmail("test@test.com")

	assert.IsType(s.T(), &common.ErrEntityNotFound{}, err)
}

func (s *UserRepositoryTestSuite) TestUserRepository_UpdateSuccessfully() {
	user := factory.NewUser()

	err := s.repo.Save(&user)
	require.Nil(s.T(), err)

	user.FirstName = "new name"
	user.LastName = "new last name"

	err = s.repo.Update(&user)
	require.Nil(s.T(), err)

	persisted, err := s.repo.FindByEmail(user.Email)
	require.Nil(s.T(), err)

	assert.Equal(s.T(), user, persisted)
}

func (s *UserRepositoryTestSuite) TestUserRepository_UpdateWhenUserDoesNotExist() {
	user := factory.NewUser()

	err := s.repo.Update(&user)
	assert.Nil(s.T(), err)
}
