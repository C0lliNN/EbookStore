package persistence_test

import (
	"context"
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/c0llinn/ebook-store/internal/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type UserRepositoryTestSuite struct {
	RepositoryTestSuite
	repo *persistence.UserRepository
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	s.RepositoryTestSuite.SetupSuite()

	s.repo = persistence.NewUserRepository(s.db)
}
func (s *UserRepositoryTestSuite) TearDownTest() {
	s.db.Delete(&auth.User{}, "1 = 1")
}

func TestUserRepositoryTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) TestUserRepository_SaveSuccessfully() {
	ctx := context.TODO()

	user := auth.User{
		ID:        "some-id",
		FirstName: "Raphael",
		LastName:  "Collin",
		Email:     "raphael@test.com",
		Role:      auth.Customer,
		Password:  "password",
		CreatedAt: time.Now().Unix(),
	}

	err := s.repo.Save(ctx, &user)

	assert.Nil(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestUserRepository_SaveWithDuplicateEmail() {
	ctx := context.TODO()

	user := auth.User{
		ID:        "some-id",
		FirstName: "Raphael",
		LastName:  "Collin",
		Email:     "raphael@test.com",
		Role:      auth.Customer,
		Password:  "password",
		CreatedAt: time.Now().Unix(),
	}

	err := s.repo.Save(ctx, &user)
	assert.Nil(s.T(), err)

	user.ID = "new-id"

	err = s.repo.Save(ctx, &user)
	assert.IsType(s.T(), &persistence.ErrDuplicateKey{}, err)
}

func (s *UserRepositoryTestSuite) TestUserRepository_FindByEmailSuccessfully() {
	ctx := context.TODO()

	expected := auth.User{
		ID:        "some-id",
		FirstName: "Raphael",
		LastName:  "Collin",
		Email:     "raphael@test.com",
		Role:      auth.Customer,
		Password:  "password",
		CreatedAt: time.Now().Unix(),
	}

	err := s.repo.Save(ctx, &expected)
	require.Nil(s.T(), err)

	actual, err := s.repo.FindByEmail(ctx, expected.Email)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)
}

func (s *UserRepositoryTestSuite) TestUserRepository_FindByEmailNotFound() {
	ctx := context.TODO()

	_, err := s.repo.FindByEmail(ctx, "test@test.com")

	assert.IsType(s.T(), &persistence.ErrEntityNotFound{}, err)
}

func (s *UserRepositoryTestSuite) TestUserRepository_UpdateSuccessfully() {
	ctx := context.TODO()

	user := auth.User{
		ID:        "some-id",
		FirstName: "Raphael",
		LastName:  "Collin",
		Email:     "raphael@test.com",
		Role:      auth.Customer,
		Password:  "password",
		CreatedAt: time.Now().Unix(),
	}

	err := s.repo.Save(ctx, &user)
	require.Nil(s.T(), err)

	user.FirstName = "new name"
	user.LastName = "new last name"

	err = s.repo.Update(ctx, &user)
	require.Nil(s.T(), err)

	persisted, err := s.repo.FindByEmail(ctx, user.Email)
	require.Nil(s.T(), err)

	assert.Equal(s.T(), user, persisted)
}

func (s *UserRepositoryTestSuite) TestUserRepository_UpdateWhenUserDoesNotExist() {
	ctx := context.TODO()

	user := auth.User{
		ID:        "some-id",
		FirstName: "Raphael",
		LastName:  "Collin",
		Email:     "raphael@test.com",
		Role:      auth.Customer,
		Password:  "password",
		CreatedAt: time.Now().Unix(),
	}

	err := s.repo.Update(ctx, &user)
	assert.Nil(s.T(), err)
}
